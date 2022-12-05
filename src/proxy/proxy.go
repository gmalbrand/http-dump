package proxy

import (
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	hopHeader = []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authorization",
		"Proxy-Authenticate",
	}
)


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


func deleteHopHeader(req *http.Request){
	for _, header := range(hopHeader){
		req.Header.Del(header)
	}
}

func addForwardHeader(req *http.Request){
	req.Header.Add("X-Forwarded-For", getRemoteAddress(req))	
}

func copyHeader(src, dst http.Header){
	for header, values := range src {
		for _, value := range values {
			dst.Add(header, value)
		}
	}
}

func ProxyHandler(w http.ResponseWriter, req *http.Request){
	log.Printf("Handling proxy request")

	deleteHopHeader(req)
	addForwardHeader(req)

	client := &http.Client{}
	// Request URI Must be removed
	req.RequestURI = ""

	resp, proxyErr := client.Do(req)

	if proxyErr != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Printf("Fatal proxy error: %s", proxyErr)
	}
	defer resp.Body.Close()

	log.Printf("Request from %s status %d", getRemoteAddress(resp.Request), resp.StatusCode)

	copyHeader(resp.Header, w.Header())
	io.Copy(w, resp.Body)
}


