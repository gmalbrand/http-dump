package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"log"

	"github.com/gmalbrand/http-dump/generators"
	"github.com/gmalbrand/http-dump/logger"
	"github.com/gmalbrand/http-dump/monitoring"
	"github.com/gmalbrand/http-dump/proxy"
)

const (
	httpDefaultPort = 8080
	defaultVersion  = "v1.0.0"
	defaultDuration = "5m"
)

var (
	version = os.Getenv("HTTP_DUMP_VERSION")
)

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

	LoadGenerator := generators.NewLoadGenerator()
	//memLoadGen := generators.NewMemLoadGenerator()

	mux := monitoring.NewMonitoredMux()

	mux.HandleFunc("/", proxy.ProxyHandler)
	mux.HandleFunc("/dump", dumpRequest)
	mux.HandleFunc("/info", info)
	mux.HandleFunc("/cpuLoad", LoadGenerator.GenerateCPULoad)
	mux.HandleFunc("/memLoad", LoadGenerator.GenerateMemLoad)
	fmt.Printf("Serving requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), logger.AccessCombinedLog(mux.Server())))
	LoadGenerator.Wait()
}
