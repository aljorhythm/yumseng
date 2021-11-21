package utils

import (
	"encoding/json"
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

func HttpResponseToStruct(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	err := json.NewDecoder(resp.Body).Decode(v)

	if err == nil {
		return err
	}

	return nil
}

func ToJson(object interface{}) ([]byte, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	message := json.RawMessage(bytes)
	return message, nil
}
