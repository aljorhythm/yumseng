package cheers

import (
	"testing"
	"time"
)

func TestService(t *testing.T) {
	service := Service{}

	t.Run("adding cheer should increase count", func(t *testing.T) {
		previousCount := len(service.cheers)

		service.addCheer(&Cheer{
			Value:    "hello1",
			DateTime: time.Now(),
		})

		wanted := previousCount + 1
		got := len(service.getCheers())

		if wanted != got {
			t.Errorf("got: %s | wanted: %s", got, wanted)
		}
	})
}
