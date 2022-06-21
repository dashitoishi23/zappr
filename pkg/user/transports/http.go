package usertransport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	util "dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(endpoints userendpoint.Set) []commonmodels.HttpServerConfig {
	var userServers []commonmodels.HttpServerConfig

	tokenHandler := httptransport.NewServer(
		endpoints.GenerateToken,
		util.DecodeHTTPGenericRequest[userendpoint.GenerateTokenRequest],
		util.EncodeHTTPGenericResponse,
	)

	loginHandler := httptransport.NewServer(
		endpoints.ValidateLogin,
		decodeValidateLoginRequest,
		encodeValidateLoginResponse,
	)

	signupuserHandler := httptransport.NewServer(
		endpoints.SignupUser,
		util.DecodeHTTPGenericRequest[userendpoint.SignupUserRequest],
		util.EncodeHTTPGenericResponse,
	)

	return append(userServers, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: tokenHandler,
		Route: "/user/token",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: loginHandler,
		Route: "/user/login",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: signupuserHandler,
		Route: "/user",
		Methods: []string{"POST"},
	})

}

func decodeValidateLoginRequest(_ context.Context, r * http.Request) (interface{}, error){
	var req userendpoint.ValidateLoginRequest

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

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