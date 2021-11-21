package integration_tests

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

func testHttpClient() *http.Client {
	return &http.Client{}
}

type HttpRequestOptions struct {
	method string
	path   string
	filter *func(r *http.Request) error
	body   io.Reader
}

func (o *HttpRequestOptions) getMethod() string {
	if o.method == "" {
		return http.MethodGet
	}
	return o.method
}

func httpRequest(options HttpRequestOptions, t *testing.T) *http.Response {
	url := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   options.path,
	}
	req, err := http.NewRequest(options.getMethod(), url.String(), options.body)

	if err != nil {
		t.Fatalf("unable to construct request error: %#v | request: %#v", err.Error(), req)
	}

	if options.filter != nil {
		filter := *options.filter
		filter(req)
	}

	client := testHttpClient()

	t.Logf("request method: %s | url: %s", req.Method, req.URL)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("error doing request %#v %#v", err.Error(), resp)
	} else {
		t.Logf("response status: %d | content length: %d", resp.StatusCode, resp.ContentLength)
	}

	return resp
}
