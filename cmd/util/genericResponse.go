package util

import (
	"context"
	"encoding/json"
	"net/http"
)

func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errResp := make(map[string]string)

		errResp["errMsg"] = e.Error()

		return json.NewEncoder(w).Encode(errResp)
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}