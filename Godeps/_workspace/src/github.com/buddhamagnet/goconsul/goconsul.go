package goconsul

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	consul Consul
	port   string
)

type Consul struct {
	Name string
}

func init() {
	port = os.Getenv("CONSUL_PORT")
	if port == "" {
		port = "8500"
	}
}

// RegisterService registers the service outlined in goconsul.json
// with the local consul agent.
func RegisterService() (err error) {
	// Open goconsul configuration file.
	config, err := os.Open("goconsul.json")
	defer config.Close()

	if err != nil {
		return err
	}

	// Read configuration file for storage in struct and
	// forward transmission to consul agent over HTTP API.
	data, err := ioutil.ReadAll(config)
	if err != nil {
		return err
	}
	// Marshal response into consul struct for easy
	// retrieval of service name for command line
	// confirmations.
	err = json.Unmarshal(data, &consul)
	if err != nil {
		return err
	}

	// Register the service with local consul agent.
	err = doRegistration(data)
	if err != nil {
		return err
	}

	// Finally, query the API to confirm this has been done.
	resp, err := queryService()

	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		response, _ := ioutil.ReadAll(resp.Body)
		log.Println("service registration complete, please check response:")
		log.Println(string(response))
	} else {
		log.Fatalf("service registration failed: %d\n", resp.StatusCode)
	}
	return nil
}
