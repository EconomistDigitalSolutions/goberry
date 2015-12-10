package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// RouteMap links RAML display names to HTTP handlers.
var RouteMap = map[string]http.HandlerFunc{
	"Root":          Root,
	"Version":       Version,
	"HealthCheck":   HealthCheck,
	"QueryOptional": QueryOptional,
	"QueryRequired": QueryRequired,
}

var ramlMime = "application/x-yaml"

// Version returns the build date and commit hash of the current build.
func Version(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(map[string]string{
		"message": fmt.Sprintf("build date: %s commit: %s", buildstamp, githash),
	})
	w.Write(json)
}

// Root is the hypermedia root endpoint of the service (GET).
func Root(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	switch strings.Contains(accept, ramlMime) {
	case true:
		w.Header().Set("Content-Type", ramlMime)
		http.ServeFile(w, r, "api.raml")
		return
	default:
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "api.html")
	}
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
