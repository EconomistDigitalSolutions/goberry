package journal

import (
	"net/http"
	"os"
	"time"

	stdlog "log"

	kitlog "github.com/go-kit/kit/log"
)

var (
	logger  kitlog.Logger
	Service string
)

func init() {
	logger = kitlog.NewJSONLogger(os.Stdout)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger))
}

func LogRequest(r *http.Request) {
	logger.Log("channel", "request", "service", Service, "method", r.Method, "url", r.URL.String(), "headers", r.Header, "ts", time.Now())
}

func LogChannel(channel string, message ...interface{}) {
	logger.Log("channel", channel, "service", Service, "message", message, "ts", time.Now())
}

func LogError(message string) {
	LogChannel("channel", message)
}

func LogInfo(message string) {
	LogChannel("information", message)
}

func LogWorker(message ...interface{}) {
	LogChannel("worker", message)
}
