package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/gmalbrand/http-dump/monitoring"
)

const (
	httpDefaultPort = 8080
)

func dumpRequest(w http.ResponseWriter, req *http.Request) {
	// Adding comment to generate a push
	var formatted, err = httputil.DumpRequest(req, true)

	if err != nil {
		fmt.Fprint(w, err)
	}
	w.Write(formatted)
}

func info(w http.ResponseWriter, req *http.Request) {
	// Printing info message (need to update the version at build time)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"api\": \"http-dump\", \"version\": 1.0.0}")
}

func main() {
	// Get parameters
	// Et ca continue encore et encore
	port, err := strconv.Atoi(os.Getenv("HTTP_SERVER_PORT"))

	if err != nil {
		port = httpDefaultPort
	}

	mux := monitoring.NewMonitoredMux()

	mux.HandleFunc("/dump", dumpRequest)
	mux.HandleFunc("/info", info)
	fmt.Printf("Serving requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux.Server()))
}
