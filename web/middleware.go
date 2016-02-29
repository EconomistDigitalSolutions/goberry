package web

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/EconomistDigitalSolutions/mt-utils/httputil"
	"github.com/EconomistDigitalSolutions/watchman/journal"
)

// JSONMiddleware writes the appropriate content type
// header for JSON output.
func JSONMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs the request method and URL string
// to the log output for every request.
func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricData.RequestCounter.Add(1)
		metricData.RequestConcurrentCounter.Add(1)

		start := time.Now()
		lw := &httputil.HTTPResponse{ResponseWriter: w, URL: r.URL.Path, RequestID: httputil.GetRequestID(r)}
		journal.LogRequestWithInfo(r, "requestId", lw.RequestID, "status", "received")
		h.ServeHTTP(lw, r)
		responseTime := time.Since(start).Seconds() * 1000

		rtmFloat, _ := strconv.ParseFloat(metricData.ResponseTimeMax.String(), 64)
		if responseTime > rtmFloat {
			metricData.ResponseTimeMax.Set(responseTime)
		}

		metricData.ResponseTimeSum.Add(responseTime)
		metricData.RequestConcurrentCounter.Add(-1)

		journal.LogRequestWithInfo(r, "responseCode", lw.Status, "responseTimeMs", responseTime, "requestId", lw.RequestID, "status", "responded")
	})
}

// RecoverMiddleware returns a function that runs a defer
// that captures runtime panics, logs the error and ensures
// the sever returns the appropriate 500 error.
func RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("gref panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
