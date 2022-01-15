package rooms

import (
	"github.com/aljorhythm/yumseng/utils/testingutils"
	"net/http/httptest"
	"testing"
)

func assertAllResponses(t *testing.T, recorder *httptest.ResponseRecorder) {
	assertable := testingutils.AssertableResponse{ResponseRecorder: recorder, T: t}
	assertable.ContentTypeIsJson()
}
