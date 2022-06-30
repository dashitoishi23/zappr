package util

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"github.com/go-kit/kit/endpoint"
)

func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(endpoint.Failer); ok && e.Failed() != nil {
		ErrorEncoder(ctx, e.Failed(), w)

		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func DecodeGenericHTTPIdentifierRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) <= 1 {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	} else if parts[2] == "" {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	}

	return parts[2], nil
}

func DecodeHTTPPagedRequest[T any](ctx context.Context, r *http.Request) (interface{}, error){
	var req T

	var pagedRequest commonmodels.PagedRequest[T]

	queries := r.URL.Query()

	if queries.Get("page") != ""{
		parsedPage, parseErr := strconv.Atoi(queries.Get("page"))
		if parseErr != nil {
			pagedRequest.Page = 1
		}else{
			pagedRequest.Page = parsedPage
		}
	} else{
		pagedRequest.Page = 1
	}

	if queries.Get("size") != ""{
		parsedSize, parseErr := strconv.Atoi(queries.Get("size"))
		if parseErr != nil {
			pagedRequest.Size = 5
		}else{
			pagedRequest.Size = parsedSize
		}
	} else{
		pagedRequest.Size = 5
	}

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	if err != nil {
		return req, err
	}

	pagedRequest.Entity = req

	return pagedRequest, nil


}

func DecodeHTTPGenericRequest[T any](ctx context.Context,  r *http.Request) (interface{}, error){
	var req T

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	if err == io.EOF {
		return req, nil
	}

	if err != nil {
		return req, err
	}

	return req, nil

}