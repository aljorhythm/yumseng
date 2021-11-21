package integration_tests

import (
	"os"
)

var (
	// host for testing
	host   = IntegrationTestHost()
	scheme = IntegrationTestScheme()
)

func IntegrationTestScheme() string {
	value := os.Getenv("SCHEME")
	if value == "" {
		return "http"
	} else {
		return value
	}
}

func IntegrationTestHost() string {
	value := os.Getenv("HOST")
	if value == "" {
		return "localhost:8080"
	} else {
		return value
	}
}
