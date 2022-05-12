package usertransport

import (
	"context"
	"encoding/json"
	"net/http"

	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(endpoints userendpoint.Set) http.Handler {
	m:= http.NewServeMux()
	m.Handle("/token", httptransport.NewServer(
		endpoints.GenerateToken,
		decodeGenerateTokenHTTPRequest,
		encodeHTTPGenericResponse,
	))

	return m
}

func decodeGenerateTokenHTTPRequest(_ context.Context, r *http.Request) (interface{}, error){
	var req userendpoint.GenerateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	return req, err
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}