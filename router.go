package goberry

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/EconomistDigitalSolutions/ramlapi"
	"github.com/EconomistDigitalSolutions/watchman/journal"
	"github.com/buddhamagnet/raml"
	"github.com/buddhamagnet/yaml"
	"github.com/gorilla/mux"
)

var (
	router     *mux.Router
	api        *raml.APIDefinition
	buildstamp string
	githash    string
)

// NewRouter creates a mux router, sets up
// a static handler and registers the dynamic
// routes and middleware handlers with the mux.
func NewRouter(ramlFile, build, hash string) *mux.Router {
	buildstamp = build
	githash = hash
	router = mux.NewRouter().StrictSlash(true)
	// Assemble middleware as required.
	assembleMiddleware()
	assembleRoutes(buildstamp, githash)
	return router
}

// assembleMiddleware sets up the middleware stack.
func assembleMiddleware() {
	http.Handle("/",
		JSONMiddleware(
			LoggingMiddleware(
				RecoverMiddleware(router))))
}

func assembleRoutes(build, hash string) {
	var err error

	// If bundling is enabled, read the RAML
	// from bundle.go
	if os.Getenv("BUNDLE_ASSETS") != "" {
		ramlFile := os.Getenv("RAMLFILE_NAME")

		if ramlFile == "" {
			ramlFile = "api.raml"
		}

		bundle, err := GoberryBundle.Open(ramlFile)
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadAll(bundle)
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal(data, &api)
		if err != nil {
			log.Fatal(err)
		}
		// Otherwise, read the file from the filesystem.
	} else {
		ramlPath := os.Getenv("RAMLFILE_PATH")

		if ramlPath == "" {
			ramlPath = "api.raml"
		}

		api, err = ramlapi.ProcessRAML(ramlPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	journal.LogChannel("raml-processor", fmt.Sprintf("processing API spec for %s", api.Title))
	journal.LogChannel("raml-processor", fmt.Sprintf("base URI at %s", api.BaseUri))
	ramlapi.Build(api, routerFunc)
}

func routerFunc(data map[string]string) {
	route := router.
		Methods(data["verb"]).
		Path(data["path"]).
		Handler(RouteMap[data["handler"]])

	if data["query"] != "" {
		if data["query_pattern"] != "" {
			route.Queries(data["query"], fmt.Sprintf("{%s:%s}", data["query"], data["query_pattern"]))
		} else {
			route.Queries(data["query"], "")
		}
	}
}
