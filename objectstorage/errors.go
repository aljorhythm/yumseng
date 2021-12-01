package objectstorage

import "errors"

var (
	ERROR_DATA_NOT_FOUND              = errors.New("ERROR_DATA_NOT_FOUND")
	ERROR_STORE_DATA_FAILED           = errors.New("ERROR_STORE_DATA_FAILED")
	ERROR_FAIL_TO_CONNECT_SPACE_STORE = errors.New("ERROR_STORE_DATA_FAILED")
)
