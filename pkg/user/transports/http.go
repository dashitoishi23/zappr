package usertransport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	util "dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(endpoints userendpoint.Set) []commonmodels.HttpServerConfig {
	var userServers []commonmodels.HttpServerConfig

	serverOptions := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(util.ErrorEncoder),
	}

	loginHandler := httptransport.NewServer(
		endpoints.ValidateLogin,
		DecodeLoginRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	signupuserHandler := httptransport.NewServer(
		endpoints.SignupUser,
		util.DecodeHTTPGenericRequest[userendpoint.SignupUserRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	return append(userServers, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: loginHandler,
		Route: "/user/login/{tenantIdentifier}",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: signupuserHandler,
		Route: "/user",
		Methods: []string{"POST"},
	})

}


func DecodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req userendpoint.ValidateLoginRequest

	parts := strings.Split(r.URL.Path, "/")

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	if len(parts) <= 1 {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	} else if parts[len(parts) - 1] == "" {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	}

	req.CurrentUser.TenantIdentifier = parts[len(parts) - 1]

	return req, err
}