package integration_tests

import (
	"os"
	"testing"
)

func IntegrationTestHost(t *testing.T) string {
	value := os.Getenv("HOST")
	if value == "" {
		t.Fatal("[main_test_utils.go] Unable to retrieve host")
		return ""
	} else {
		return value
	}
}
