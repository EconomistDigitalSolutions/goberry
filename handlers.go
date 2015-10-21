package goberry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var RouteMap = map[string]http.HandlerFunc{

	"Root":    Root,
	"Version": Version,
}

// Handler for rest URI /version and the action GET
func Version(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(map[string]string{
		"message": fmt.Sprintf("build date: %s commit: %s", buildstamp, githash),
	})
	w.Write(json)
}

// Handler for rest URI / and the action GET
func Root(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(map[string]string{
		"message": "RootGET",
	})
	w.Write(json)
}
