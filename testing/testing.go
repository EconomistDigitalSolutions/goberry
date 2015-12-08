package testing

import (
	"io/ioutil"
	"log"
	"testing"
)

// WithContext takes a slice of functions and a test
// function and sets up the context in which a test runs
// allowing setup and teardown operations.
func WithContext(t *testing.T, setup []func(), teardown []func(), test func(*testing.T)) {
	runFuncs(setup)
	test(t)
	runFuncs(teardown)
}

// DisableLogger redirects the logger to subspace.
func DisableLogger() {
	log.SetOutput(ioutil.Discard)
}

func runFuncs(funcs []func()) {
	for _, fn := range funcs {
		fn()
	}
}
