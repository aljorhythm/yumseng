package integration_tests

import (
	"bytes"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCheer(t *testing.T) {
	t.Run("GET /cheers should return a list of Cheers", func(t *testing.T) {
		response := httpRequest(HttpRequestOptions{
			path: "/cheers",
		}, t)

		got := []*cheers.Cheer{}
		err := utils.HttpResponseBodyToStruct(response, &got)

		assert.NoError(t, err)

		wanted := []*cheers.Cheer{}

		assert.Equal(t, wanted, got)

		t.Run("POST /cheers with valid request should be good", func(t *testing.T) {
			cheer := cheers.Cheer{
				Value: "yum",
			}
			message, err := utils.ToJson(cheer)

			assert.NoError(t, err)

			response := httpRequest(HttpRequestOptions{
				path:   "/cheers",
				method: http.MethodPost,
				body:   bytes.NewBuffer(message),
			}, t)

			assert.Equal(t, http.StatusOK, response.StatusCode)

			t.Run("GET /cheers should not return a list with length 1", func(t *testing.T) {

				response := httpRequest(HttpRequestOptions{
					path:   "/cheers",
					method: http.MethodGet,
				}, t)

				err = utils.HttpResponseBodyToStruct(response, &got)
				assert.NoError(t, err)

				assert.Equal(t, http.StatusOK, response.StatusCode)
				assert.Equal(t, 1, len(got))
			})
		})
	})
}
