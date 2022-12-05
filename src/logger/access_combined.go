package logger

import (
	"log"
	"net/http"
	"strings"
	"time"
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
