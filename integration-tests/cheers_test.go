package integration_tests

import (
	"github.com/aljorhythm/yumseng/utils"
	"testing"
)

func TestCheer(t *testing.T) {
	t.Run("GET /cheers should return 'Cheers'", func(t *testing.T) {
		response := httpRequest(HttpRequestOptions{
			path: "/cheers",
		}, t)

		got, err := utils.HttpResponseToString(response)

		if err != nil {
			t.Fatal()
		}

		wanted := "Cheers!"

		if got != wanted {
			t.Errorf("got: %s | wanted: %s", got, wanted)
		} else {
			t.Logf("nice cheers!")
		}
	})
}
