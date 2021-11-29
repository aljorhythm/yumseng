package cheers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	service := service{}

	t.Run("adding cheer should increase count", func(t *testing.T) {
		previousCount := len(service.cheers)

		service.AddCheer(&Cheer{
			Value:           "hello1",
			ClientCreatedAt: time.Now(),
		})

		wanted := previousCount + 1
		got := len(service.GetCheers())

		assert.Equal(t, wanted, got)
	})
}
