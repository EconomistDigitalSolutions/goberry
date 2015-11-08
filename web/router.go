package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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
	assembleRoutes()
	return router
}

// assembleMiddleware sets up the middleware stack.
func assembleMiddleware() {
	http.Handle("/",
		JSONMiddleware(
			LoggingMiddleware(
				RecoverMiddleware(router))))
}

func assembleRoutes() {
	var err error

	// If bundling is enabled, read the RAML
	// from bundle.go
	if os.Getenv("BUNDLE_ASSETS") == "1" {
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

		api, err = ramlapi.Process(ramlPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	journal.LogChannel("raml-processor", fmt.Sprintf("processing API spec for %s", api.Title))
	journal.LogChannel("raml-processor", fmt.Sprintf("base URI at %s", api.BaseUri))
	ramlapi.Build(api, routerFunc)
}

func routerFunc(ep *ramlapi.Endpoint) {

	path := ep.Path

	for _, up := range ep.URIParameters {
		if up.Pattern != "" {
			path = strings.Replace(
				path,
				fmt.Sprintf("{%s}", up.Key),
				fmt.Sprintf("{%s:%s}", up.Key, up.Pattern),
				1)
		}
	}

	route := router.
		Methods(ep.Verb).
		Path(path).
		Handler(RouteMap[ep.Handler])

	for _, param := range ep.QueryParameters {
		if param.Pattern != "" {
			route.Queries(param.Key, fmt.Sprintf("{%s:%s}", param.Key, param.Pattern))
		} else {
			route.Queries(param.Key, "")
		}
	}
}
