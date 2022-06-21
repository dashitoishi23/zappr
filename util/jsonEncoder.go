package util

import (
	"bytes"
	"encoding/json"
)

func JsonEncoder[T any](obj T) ([]byte, error) {
	encodedBytesBuffer := new(bytes.Buffer)

	if err := json.NewEncoder(encodedBytesBuffer).Encode(obj); err != nil {
		return encodedBytesBuffer.Bytes(), err
	}

	return encodedBytesBuffer.Bytes(), nil

}