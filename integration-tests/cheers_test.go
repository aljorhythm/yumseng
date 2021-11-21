package integration_tests

import (
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheer(t *testing.T) {
	t.Run("GET /cheers should return a list of Cheers", func(t *testing.T) {
		response := httpRequest(HttpRequestOptions{
			path: "/cheers",
		}, t)

		got := []*cheers.Cheer{}
		err := utils.HttpResponseToStruct(response, &got)

		if err != nil {
			t.Error(err)
		}

		wanted := []*cheers.Cheer{}

		assert.Equal(t, wanted, got)
	})
}
