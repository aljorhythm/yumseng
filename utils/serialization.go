package utils

import (
	"encoding/json"
	"fmt"
)

func MustEncodeJson(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("MustEncodeJson error %#v", v))
	}
	return bytes
}
