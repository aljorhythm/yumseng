package integration_tests

import "testing"

func TestGetHost(t *testing.T) {
	host := IntegrationTestHost()

	if host != "" {
		t.Logf("Able to get host from environment for testing: %s", host)
	} else {
		t.Errorf("Unable to get host from environment for testing")
	}
}
