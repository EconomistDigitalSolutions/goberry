package goberry_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/EconomistDigitalSolutions/goberry"
	"github.com/gorilla/mux"
)

type Response struct {
	Message string
}

var (
	router   *mux.Router
	response Response
)

func TestMain(m *testing.M) {
	router = goberry.NewRouter("api.raml", "", "")
	retCode := m.Run()
	os.Exit(retCode)
}

func TestRootEndpoint(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	json.NewDecoder(res.Body).Decode(&response)
	if response.Message != "RootGET" {
		t.Fatalf("expected RootGET as message, got %s\n", response.Message)
	}
}
