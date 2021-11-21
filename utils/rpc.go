package utils

import (
	"io"
	"net/http"
)

/**
response will be read and closed
*/
func HttpResponseToString(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}
