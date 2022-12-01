package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gmalbrand/http-dump/monitoring"
)

const (
	httpDefaultPort = 8080
	defaultVersion = "v1.0.0"
)

var (
	version = os.Getenv("HTTP_DUMP_VERSION")
)

type (
	loggingResponseData struct{
		size int
		status int
	}

	loggingResponseWriter struct{
		http.ResponseWriter
		loggingData *loggingResponseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
    size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
    r.loggingData.size += size // capture size
    return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
    r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
    r.loggingData.status = statusCode // capture status code
}

func dumpRequest(w http.ResponseWriter, req *http.Request) {
	// Adding comment to generate a push and another one
	var formatted, err = httputil.DumpRequest(req, true)

	if err != nil {
		fmt.Fprint(w, err)
	}
	w.Write(formatted)
}

func info(w http.ResponseWriter, req *http.Request) {
	// Printing info message (need to update the version at build time)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"api\": \"http-dump\", \"version\": \"%s\"}", version)
}

func getDefault(value string) string{
	if value == "" {
		return "-"
	}else{
		return value
	}
}

func getRemoteAddress(r *http.Request) string{
	hostAddress := "-"
		addressSlice := strings.Split(r.RemoteAddr, ":")
		if len(addressSlice) > 2 {
			hostAddress = strings.Join(addressSlice[:len(addressSlice)-1], ":")
		}else{
			hostAddress = addressSlice[0]
		}
	return hostAddress
}

func AccessCombinedLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time := time.Now().Format("02/Jan/2006:15:04:05 -0700")
		loggingData := &loggingResponseData{size : 0, status: 0}
		lw := loggingResponseWriter{ResponseWriter: w, loggingData: loggingData}
		handler.ServeHTTP(&lw, r)

		referer := getDefault(r.Referer())
		
		user,_,_ := r.BasicAuth()
		user = getDefault(user)

		log.Printf("%s - %s [%s] \"%s  %s %s\" %d %d \"%s\" \"%s\"", getRemoteAddress(r), user, time, r.Method, r.URL, r.Proto, loggingData.status, loggingData.size, referer, r.UserAgent())
    })
}


func main() {
	// Get parameters
	port, err := strconv.Atoi(os.Getenv("HTTP_SERVER_PORT"))
	log.SetFlags(0)
	
	if err != nil {
		port = httpDefaultPort
	}

	if version == "" {
		version = defaultVersion
	}

	mux := monitoring.NewMonitoredMux()

	mux.HandleFunc("/dump", dumpRequest)
	mux.HandleFunc("/info", info)
	fmt.Printf("Serving requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), AccessCombinedLog(mux.Server())))
}
