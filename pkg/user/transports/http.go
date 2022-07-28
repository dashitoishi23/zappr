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
		util.DecodeHTTPGenericRequest[userendpoint.RegisterOAuthRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	authenticateGoogleOAuthRedirect := httptransport.NewServer(
		endpoints.AuthenticateGoogleOAuthRedirect,
		DecodeGoogleOAuthRedirect,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	authenticateGoogleAccessToken := httptransport.NewServer(
		endpoints.AuthenticateGoogleAccessToken,
		DecodeAuthenticateGoogleAccessTokenRequest,
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
		Route: "/oauth",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: authenticateGoogleOAuthRedirect,
		Route: "/oauth/google/callback",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: authenticateGoogleAccessToken,
		Route: "/oauth/google/accesstoken/{tenantIdentifier}",
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

func DecodeGoogleOAuthRedirect(ctx context.Context, r *http.Request) (interface{}, error) {
	var req userendpoint.AuthenticateGoogleOAuthRedirectRequest

	queries := r.URL.Query()

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

func DecodeAuthenticateGoogleAccessTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req userendpoint.AuthenticateGoogleAccessTokenRequest

	parts := strings.Split(r.URL.Path, "/")

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	if len(parts) < 3 {
		return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
	}

	req.TenantIdentifier = parts[len(parts) - 1]

	return req, err

}