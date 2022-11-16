package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MonitoredMux struct {
	server          *http.ServeMux
	inFlightRequest prometheus.Gauge
	requestCount    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
}

func NewMonitoredMux() *MonitoredMux {
	m := new(MonitoredMux)

	m.server = http.NewServeMux()
	m.inFlightRequest = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "http_requests_in_flight",
		Help:        "A gauge of requests currently being served by the wrapped handler.",
		ConstLabels: prometheus.Labels{"server": "api"},
	})
	prometheus.MustRegister(m.inFlightRequest)

	m.requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "http_requests_total",
			Help:        "A counter for requests to the wrapped handler.",
			ConstLabels: prometheus.Labels{"server": "api"},
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(m.requestCount)

	m.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "http_request_duration_seconds",
			Help:        "A histogram of latencies for requests.",
			Buckets:     []float64{.25, .5, 1, 2.5, 5, 10},
			ConstLabels: prometheus.Labels{"server": "api"},
		},
		[]string{"handler", "method"},
	)

	prometheus.MustRegister(m.requestDuration)

	m.responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "http_response_size_bytes",
			Help:        "A histogram of response sizes for requests.",
			Buckets:     []float64{200, 500, 900, 1500},
			ConstLabels: prometheus.Labels{"server": "api"},
		},
		[]string{},
	)

	prometheus.MustRegister(m.responseSize)

	m.requestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "http_request_size_bytes",
			Help:        "A histogram of requests sizes.",
			Buckets:     []float64{200, 500, 900, 1500},
			ConstLabels: prometheus.Labels{"server": "api"},
		},
		[]string{},
	)

	prometheus.MustRegister(m.requestSize)

	m.server.Handle("/metrics", promhttp.Handler())
	return m
}

func (m *MonitoredMux) HandleFunc(pattern string, handlerFunc func(arg1 http.ResponseWriter, arg2 *http.Request)) {
	handlerChain := promhttp.InstrumentHandlerInFlight(m.inFlightRequest,
		promhttp.InstrumentHandlerDuration(m.requestDuration.MustCurryWith(prometheus.Labels{"handler": pattern[1:]}),
			promhttp.InstrumentHandlerCounter(m.requestCount,
				promhttp.InstrumentHandlerResponseSize(m.responseSize,
					promhttp.InstrumentHandlerRequestSize(m.requestSize, http.HandlerFunc(handlerFunc)),
				),
			),
		),
	)

	m.server.Handle(pattern, handlerChain)
}

func (m *MonitoredMux) Server() *http.ServeMux {
	return m.server
}
