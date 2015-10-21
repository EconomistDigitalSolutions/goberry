package goberry

import (
	"log"
	"net/http"

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
		h.ServeHTTP(w, r)
		journal.LogRequest(r)
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
