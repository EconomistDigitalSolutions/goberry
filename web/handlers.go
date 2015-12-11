package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RouteMap links RAML display names to HTTP handlers.
var RouteMap = map[string]http.HandlerFunc{
	"Root":          Root,
	"Version":       Version,
	"HealthCheck":   HealthCheck,
	"QueryOptional": QueryOptional,
	"QueryRequired": QueryRequired,
}

// Version returns the build date and commit hash of the current build.
func Version(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(map[string]string{
		"message": fmt.Sprintf("build date: %s commit: %s", buildstamp, githash),
	})
	w.Write(json)
}

// Root is the hypermedia root endpoint of the service (GET).
func Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "api.html")
	return
}

// HealthCheck is the health check endpoint of the service (HEAD).
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	up := true
	if up {
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

// QueryOptional demonstrates an endpoint with optional query parameters.
func QueryOptional(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(map[string]string{
		"message": "query optional",
	})
	w.Write(json)
}

// QueryRequired demonstrates an endpoint with required query parameters.
func QueryRequired(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(map[string]string{
		"message": "query required",
	})
	w.Write(json)
}
