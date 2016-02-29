package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/EconomistDigitalSolutions/mt-utils/metricutil"
	"github.com/EconomistDigitalSolutions/ramlapi"
	"github.com/EconomistDigitalSolutions/watchman/journal"
	"github.com/buddhamagnet/raml"
	"github.com/buddhamagnet/yaml"
	"github.com/gorilla/mux"
)

var (
	// Router is the package level mux router.
	Router     *mux.Router
	api        *raml.APIDefinition
	buildstamp string
	githash    string
	metricData *metricutil.MetricData
)

// NewRouter creates a mux router, sets up
// a static handler and registers the dynamic
// routes and middleware handlers with the mux.
func NewRouter(ramlFile, build, hash string) {
	buildstamp = build
	githash = hash
	Router = mux.NewRouter().StrictSlash(true)
	// Create and configure metrics.
	metricData = metricutil.NewMetricData()
	metricutil.InitMetricData(metricData)
	metricData.ServiceName = journal.Service
	// Assemble middleware as required.
	assembleMiddleware()
	assembleRoutes()
}

// assembleMiddleware sets up the middleware stack.
func assembleMiddleware() {
	http.Handle("/",
		JSONMiddleware(
			LoggingMiddleware(
				RecoverMiddleware(Router))))
	// Kick off metrics.
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				metricData.SendMetricDataToCloudWatch(metricData)
				// Update http request metrics to zero everytime you send to cloudwatch to track maximum for each interval
				metricutil.InitMetricData(metricData)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
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
	err = ramlapi.Build(api, routerFunc)
	if err != nil {
		log.Fatal(err)
	}
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

	route := Router.
		Methods(ep.Verb).
		Path(path).
		Handler(RouteMap[ep.Handler])

	for _, param := range ep.QueryParameters {
		if param.Required {
			if param.Pattern != "" {
				route.Queries(param.Key, fmt.Sprintf("{%s:%s}", param.Key, param.Pattern))
			} else {
				route.Queries(param.Key, "")
			}
		}
	}
}
