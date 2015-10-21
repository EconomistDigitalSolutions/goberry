package goconsul

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

var (
	endpointRegister string
	endpointQuery    string
	client           *http.Client
)

// NewClient returns a client struct with
// a preset timeout.
func NewClient() *http.Client {
	return &http.Client{}
}

// doRegistration registers a service
// with the local consul agent.
func doRegistration(data []byte) (err error) {
	// Register service via PUT request.
	endpointRegister = fmt.Sprintf("http://localhost:%s/v1/agent/service/register", port)
	r, err := http.NewRequest("PUT", endpointRegister, bytes.NewBufferString(string(data)))

	if err != nil {
		return err
	}

	log.Printf("sending service registration request for %s\n", consul.Name)

	client = NewClient()
	_, err = client.Do(r)

	if err != nil {
		return err
	}
	return nil
}

// queryService queries the local consul agent for a service.
func queryService() (resp *http.Response, err error) {
	// Query service via HTTP API for confirmation in console.
	endpointQuery = fmt.Sprintf("http://localhost:%s/v1/catalog/service/%s", port, consul.Name)
	req, err := http.NewRequest("GET", endpointQuery, nil)
	if err != nil {
		return nil, err
	}
	client = NewClient()
	resp, err = client.Do(req)

	if err != nil {
		return nil, err
	}
	return resp, nil
}
