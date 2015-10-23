package goberry_test

import (
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
	if res.Code != http.StatusOK {
		t.Fatalf("expected 200 response got %d\n", res.Code)
	}
}
