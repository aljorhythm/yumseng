package utils

import (
	"bytes"
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

func HttpResponseBodyToStruct(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	err := DecodeJson(resp.Body, v)

	if err != nil {
		return err
	}

	return nil
}

func HttpRequestBodyToStruct(resp *http.Request, v interface{}) error {
	defer resp.Body.Close()

	err := DecodeJson(resp.Body, v)

	if err != nil {
		return err
	}

	return nil
}

func DecodeJson(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func DecodeJsonFromBytes(data []byte, v interface{}) error {
	r := bytes.NewReader(data)
	return json.NewDecoder(r).Decode(v)
}

func ToJson(object interface{}) ([]byte, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	message := json.RawMessage(bytes)
	return message, nil
}
