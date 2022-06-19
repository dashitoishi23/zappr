package util

import (
	"bytes"
	"encoding/json"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
)

func JsonEncoder[T any](obj *commonmodels.Entity[T]) ([]byte, error) {
	encodedBytesBuffer := new(bytes.Buffer)

	if err := json.NewEncoder(encodedBytesBuffer).Encode(obj); err != nil {
		return encodedBytesBuffer.Bytes(), err
	}

	return encodedBytesBuffer.Bytes(), nil

}