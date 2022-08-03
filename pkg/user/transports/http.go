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

	updateUserRoleHandler := httptransport.NewServer(
		endpoints.UpdateUserRole,
		util.DecodeHTTPGenericRequest[userendpoint.UpdateUserRoleRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	generateAPIKeyHandler := httptransport.NewServer(
		endpoints.GenerateAPIKey,
		util.DecodeHTTPGenericRequest[userendpoint.GenerateAPIKeyRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	loginWithAPIKeyHandler := httptransport.NewServer(
		endpoints.LoginWithAPIKey,
		util.DecodeHTTPGenericRequest[userendpoint.LoginWithAPIKeyRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	registerOAuth := httptransport.NewServer(
		endpoints.RegisterOAuth,
		DecodeRegisterOAuthRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	authenticateOAuthRedirect := httptransport.NewServer(
		endpoints.AuthenticateOAuthRedirect,
		DecodeOAuthRedirect,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	authenticateAccessToken := httptransport.NewServer(
		endpoints.AuthenticateAccessToken,
		DecodeAuthenticateAccessTokenRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	updateUser := httptransport.NewServer(
		endpoints.UpdateUser,
		util.DecodeHTTPGenericRequest[userendpoint.UpdateUserRequest],
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
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: updateUserRoleHandler,
		Route:"/user/role",
		Methods: []string{"PUT"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: generateAPIKeyHandler,
		Route: "/user/apikey",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: loginWithAPIKeyHandler,
		Route: "/user/apikey/login",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: registerOAuth,
		Route: "/oauth/{tenantIdentifier}",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: authenticateOAuthRedirect,
		Route: "/oauth/{providerName}/callback/{tenantIdentifier}",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: authenticateAccessToken,
		Route: "/oauth/{providerName}/accesstoken/{tenantIdentifier}",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: updateUser,
		Route: "/user",
		Methods: []string{"PUT"},
	})

}


func DecodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req userendpoint.ValidateLoginRequest

	parts := strings.Split(r.URL.Path, "/")

	defer r.Body.Close()

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

func DecodeOAuthRedirect(ctx context.Context, r *http.Request) (interface{}, error) {
	var req userendpoint.AuthenticateGoogleOAuthRedirectRequest

	queries := r.URL.Query()

	parts := strings.Split(r.URL.Path, "/")

	req.ProviderName = parts[len(parts) - 3]
	req.TenantIdentifier = parts[len(parts) - 1]
	req.Host = r.Host

	uriCode := queries.Get("code")
	uriState := queries.Get("state")

	if uriCode != "" {
		req.Code = uriCode
	} else {
		return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
	}

	if uriState != "" {
		req.State = uriState
	} else {
		return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
	}

	return req, nil
}

func DecodeAuthenticateAccessTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req userendpoint.AuthenticateAccessTokenRequest

	parts := strings.Split(r.URL.Path, "/")

	defer r.Body.Close()

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	if len(parts) < 3 {
		return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
	}

	req.TenantIdentifier = parts[len(parts) - 1]
	req.ProviderName = parts[len(parts) - 3]

	return req, err

}

func DecodeRegisterOAuthRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req userendpoint.RegisterOAuthRequest

	parts := strings.Split(r.URL.Path, "/")

	req.TenantIdentifier = parts[len(parts) - 1]

	req.Host = r.Host

	defer r.Body.Close()

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	return req, err

}