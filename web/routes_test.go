package web_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/EconomistDigitalSolutions/goberry/web"

	. "github.com/EconomistDigitalSolutions/goberry/testing"
)

var res *httptest.ResponseRecorder

// RouterSetup boots the router and fires up
// an HTTP response recorder for each test.
func routerSetup() {
	web.NewRouter("", "", "")
}

func responseSetup() {
	res = httptest.NewRecorder()
}

func setRAMLPath() {
	os.Setenv("RAMLFILE_PATH", "../api.raml")
}

func resetRAMLPath() {
	os.Setenv("RAMLFILE_PATH", "")
}

func TestRootHandler(t *testing.T) {
	WithContext(t, []func(){DisableLogger, setRAMLPath, routerSetup, responseSetup}, []func(){resetRAMLPath}, func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		// Call the endpoint.
		web.Router.ServeHTTP(res, req)
		// Check status code from root endpoint.
		if res.Code != http.StatusOK {
			t.Errorf("expected root endpoint (GET) to return 200, got %d", res.Code)
		}
	})
}

func TestVersionHandler(t *testing.T) {
	WithContext(t, []func(){setRAMLPath, responseSetup}, []func(){resetRAMLPath}, func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/version", nil)
		// Call the endpoint.
		web.Router.ServeHTTP(res, req)
		// Check status code from root endpoint.
		if res.Code != http.StatusOK {
			t.Errorf("expected version endpoint (GET) to return 200, got %d", res.Code)
		}
	})
}

func TestHealthcheckHandler(t *testing.T) {
	WithContext(t, []func(){setRAMLPath, responseSetup}, []func(){resetRAMLPath}, func(t *testing.T) {
		req, _ := http.NewRequest("HEAD", "/up", nil)
		// Call the endpoint.
		web.Router.ServeHTTP(res, req)
		// Check status code from root endpoint.
		if res.Code != http.StatusOK {
			t.Errorf("expected up endpoint (HEAD) to return 200, got %d", res.Code)
		}
	})
}

func TestInvalidHandler(t *testing.T) {
	WithContext(t, []func(){setRAMLPath, responseSetup}, []func(){resetRAMLPath}, func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/foobar", nil)
		// Call the endpoint.
		web.Router.ServeHTTP(res, req)
		// Check status code from root endpoint.
		if res.Code != http.StatusNotFound {
			t.Errorf("expected invalid endpoint to return 404, got %d", res.Code)
		}
	})
}

func TestInvalidVerbHandler(t *testing.T) {
	WithContext(t, []func(){setRAMLPath, responseSetup}, []func(){resetRAMLPath}, func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", nil)
		// Call the endpoint.
		web.Router.ServeHTTP(res, req)
		// Check status code from root endpoint.
		if res.Code != http.StatusNotFound {
			t.Errorf("expected invalid verb (POST /) to return 404, got %d", res.Code)
		}
	})
}
