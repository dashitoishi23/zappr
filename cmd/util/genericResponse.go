package util

import (
	"context"
	"encoding/json"
	"net/http"

	constants "dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/constants"
)

func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		switch(e.Error()){
		case constants.INVALID_MODEL:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		errResp := make(map[string]string)

		errResp["errMsg"] = e.Error()

		return json.NewEncoder(w).Encode(errResp)
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}