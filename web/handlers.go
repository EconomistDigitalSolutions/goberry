package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var RouteMap = map[string]httprouter.Handle{

	"Root":          Root,
	"Version":       Version,
	"HealthCheck":   HealthCheck,
	"QueryOptional": QueryOptional,
	"QueryRequired": QueryRequired,
}

// Handler for rest URI /version and the action GET
func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json, _ := json.Marshal(map[string]string{
		"message": fmt.Sprintf("build date: %s commit: %s", buildstamp, githash),
	})
	w.Write(json)
}

// Handler for rest URI / and the action GET
func Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "api.html")
	return
}

// Handler for rest URI /healthcheck and the action HEAD
func HealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	up := true
	if up {
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func QueryOptional(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json, _ := json.Marshal(map[string]string{
		"message": "query optional",
	})
	w.Write(json)
}

func QueryRequired(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json, _ := json.Marshal(map[string]string{
		"message": "query required",
	})
	w.Write(json)
}
