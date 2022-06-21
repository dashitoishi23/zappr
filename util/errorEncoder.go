package util

import (
	"context"
	"encoding/json"
	"net/http"

	constants "dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
)

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(setErrorResponseHeader(err.Error()))

	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func setErrorResponseHeader(msg string) int {
	switch(msg){
	case constants.INVALID_MODEL:
		return http.StatusBadRequest
	case constants.RECORD_NOT_FOUND:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

type errorWrapper struct {
	Error string `json:"error"`
}