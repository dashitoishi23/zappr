package usertransport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(endpoints userendpoint.Set) http.Handler {
	m:= http.NewServeMux()

	m.Handle("/token", httptransport.NewServer(
		endpoints.GenerateToken,
		decodeGenerateTokenHTTPRequest,
		encodeHTTPGenericResponse,
	))

	m.Handle("/login", httptransport.NewServer(
		endpoints.ValidateLogin,
		decodeValidateLoginRequest,
		encodeValidateLoginResponse,
	))

	return m
}

func decodeGenerateTokenHTTPRequest(_ context.Context, r *http.Request) (interface{}, error){
	var req userendpoint.GenerateTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err == io.EOF{
		return req, nil	
	}
	
	return req, err
}

func decodeValidateLoginRequest(_ context.Context, r * http.Request) (interface{}, error){
	var req userendpoint.ValidateLoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	token := r.Header.Get("Authorization")

	if token != "" {
		req.Token = token

		return req, nil
	}

	if err == io.EOF{
		return req, nil	
	}

	return req, err
}

func encodeValidateLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(userendpoint.ValidateLoginResponse)
	
	if !resp.IsValid {
		w.WriteHeader(http.StatusUnauthorized)

		return json.NewEncoder(w).Encode(response)
	}
	
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}