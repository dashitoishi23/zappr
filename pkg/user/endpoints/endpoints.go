package userendpoint

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	usermodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	util "dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	ValidateLogin endpoint.Endpoint
	SignupUser endpoint.Endpoint
	UpdateUserRole endpoint.Endpoint
	GenerateAPIKey endpoint.Endpoint
	LoginWithAPIKey endpoint.Endpoint
	RegisterOAuth endpoint.Endpoint
	AuthenticateGoogleOAuthRedirect endpoint.Endpoint
	AuthenticateGoogleAccessToken endpoint.Endpoint
} //defines all endpoints as having type Endpoint, provided by go-kit

func New(svc userservice.UserService, logger log.Logger) Set {

	var validateLoginEndpoint endpoint.Endpoint
	{
		validateLoginEndpoint = ValidateLoginEndpoint(svc)

		validateLoginEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "validateToken"))(validateLoginEndpoint)
	}

	var signupUserEndpoint endpoint.Endpoint
	{
		signupUserEndpoint = SignupUserEndpoint(svc)

		signupUserEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "signupUser"))(signupUserEndpoint)
	}

	var updateUserRoleEndpoint endpoint.Endpoint
	{
		updateUserRoleEndpoint = UpdateUserRoleEndpoint(svc)

		updateUserRoleEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "updateUserRole"))(updateUserRoleEndpoint)
	}

	var generateAPIKeyEndpoint  endpoint.Endpoint
	{
		generateAPIKeyEndpoint = GenerateAPIKeyEndpoint(svc)

		generateAPIKeyEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "generateAPIKey"))(generateAPIKeyEndpoint)
	}

	var loginWithAPIKeyEndpoint endpoint.Endpoint
	{
		loginWithAPIKeyEndpoint = LoginWithAPIKeyEndpoint(svc)

		loginWithAPIKeyEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "loginWithAPIKey"))(loginWithAPIKeyEndpoint)
	}

	var registerOAuthEndpoint endpoint.Endpoint
	{
		registerOAuthEndpoint = RegisterOAuthEndpoint(svc)

		registerOAuthEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "registerGoogleOauth"))(registerOAuthEndpoint)
	}

	var authenticateGoogleOAuthRedirectEndpoint endpoint.Endpoint
	{
		authenticateGoogleOAuthRedirectEndpoint = AuthenticateGoogleOAuthRedirectEndpoint(svc)

		authenticateGoogleOAuthRedirectEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "authenticateGoogleOauthRedirect"))(authenticateGoogleOAuthRedirectEndpoint)
	}

	var authenticateGoogleAccessTokenEndpoint endpoint.Endpoint
	{
		authenticateGoogleAccessTokenEndpoint = AuthenticateGoogleAccessTokenEndpoint(svc)

		authenticateGoogleAccessTokenEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "authenticateGoogleAccessToken"))(authenticateGoogleAccessTokenEndpoint)
	}

	return Set{
		ValidateLogin: validateLoginEndpoint,
		SignupUser: signupUserEndpoint,
		UpdateUserRole: updateUserRoleEndpoint,
		GenerateAPIKey: generateAPIKeyEndpoint,
		LoginWithAPIKey: loginWithAPIKeyEndpoint,
		RegisterOAuth: registerOAuthEndpoint,
		AuthenticateGoogleOAuthRedirect: authenticateGoogleOAuthRedirectEndpoint,
		AuthenticateGoogleAccessToken: authenticateGoogleAccessTokenEndpoint,
	}
}

func ValidateLoginEndpoint(s userservice.UserService) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ValidateLoginRequest)
		resp, err := s.LoginUser(ctx, req.CurrentUser)

		return ValidateLoginResponse{resp, err}, err
	}
}

func SignupUserEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SignupUserRequest)
		req.NewUser.InitFields()
		resp, err := s.SignupUser(ctx, req.NewUser)

		return SignupUserResponse{resp, err}, err
	}
}

func UpdateUserRoleEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req:= request.(UpdateUserRoleRequest)

		if err != nil {
			return nil, err
		}

		updatedUser, updateErr := s.UpdateUserRole(ctx, req.MasterRoleIdentifier, req.UserIdentifier)

		return UpdateUserRoleResponse{updatedUser, updateErr}, updateErr
	}
}

func GenerateAPIKeyEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GenerateAPIKeyRequest)

		apiKey, err := s.GenerateAPIKey(ctx, req.APIKeyName)

		return GenerateAPIKeyResponse{apiKey, err}, err
	}
}

func LoginWithAPIKeyEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(LoginWithAPIKeyRequest)

		jwt, err := s.LoginWithAPIKey(ctx, req.APIKey)

		return LoginWithAPIKeyResponse{jwt, err}, err
	}
}

func RegisterOAuthEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RegisterOAuthRequest)

		metadata := map[string]interface{}{
			"client_id": req.ClientID,
			"client_secret": req.ClientSecret,
			"redirect_uri": req.RedirectURI,
		}

		metadataBytes, jsonErr := json.Marshal(metadata)

		if jsonErr != nil {
			return nil, jsonErr
		}

		newOAuthProvider := usermodels.OAuthProvider{
			Name: req.ProviderName,
			Metadata: metadataBytes,
		}

		newOAuthProvider.InitFields()

		url, err := s.RegisterOAuth(ctx, newOAuthProvider, req.Scopes)

		return RegisterGoogleOAuthResponse{url, err}, err

	}
}

func AuthenticateGoogleOAuthRedirectEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthenticateGoogleOAuthRedirectRequest)

		if req.State != os.Getenv("STATE") {
			return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
		}

		jwt, err := s.AuthenticateGoogleOAuthRedirect(ctx, req.Code)

		return AuthenticateGoogleOAuthRedirectResponse{jwt, err}, err
	}

}

func AuthenticateGoogleAccessTokenEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthenticateGoogleAccessTokenRequest)

		jwt, err := s.AuthenticateGoogleAccessToken(ctx, req.AccessToken, req.TenantIdentifier)

		return ValidateLoginResponse{jwt, err}, err
	}
}

