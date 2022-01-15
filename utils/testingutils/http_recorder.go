package testingutils

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

type AssertableResponse struct {
	*httptest.ResponseRecorder
	T *testing.T
}

func (r AssertableResponse) StatusCodeIs(status string) {
	assert.Equal(r.T, status, r.Code)
}

func (r AssertableResponse) ContentTypeIsJson() {
	assert.Equal(r.T, "application/json", r.Header().Get("content-type"), "content-type should be json")
}
