package usertransport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	authhandler "dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/util"
	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)
func NewHttpHandler(endpoints userendpoint.Set) http.Handler {

	tokenHandler := httptransport.NewServer(
		endpoints.GenerateToken,
		decodeGenerateTokenHTTPRequest,
		encodeHTTPGenericResponse,
	)

	loginHandler := httptransport.NewServer(
		endpoints.ValidateLogin,
		decodeValidateLoginRequest,
		encodeValidateLoginResponse,
	)

	r := mux.NewRouter()

	r.Handle("/token", tokenHandler).Methods("POST")

	r.Handle("/login", loginHandler).Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(authhandler.AuthHandler),
		negroni.Wrap(tokenHandler),
	))

	return r
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