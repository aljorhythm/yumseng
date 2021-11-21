package integration_tests

import (
	"github.com/aljorhythm/yumseng/ping"
	"github.com/aljorhythm/yumseng/utils"
	"testing"
)

func TestGetHost(t *testing.T) {
	host := IntegrationTestHost()

	if host != "" {
		t.Logf("Able to get host from environment for testing: %s", host)
	} else {
		t.Errorf("Unable to get host from environment for testing")
	}
}

func TestPing(t *testing.T) {
	response := httpRequest(HttpRequestOptions{
		path: "/ping",
	}, t)

	got := &ping.PingResponse{}
	utils.HttpResponseBodyToStruct(response, got)

	if got.Tag == "" {
		t.Errorf("tag is empty")
	} else {
		t.Logf("tag is %s", got.Tag)
	}
}
