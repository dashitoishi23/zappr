package userendpoint

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
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
	AuthenticateOAuthRedirect endpoint.Endpoint
	AuthenticateAccessToken endpoint.Endpoint
	UpdateUser endpoint.Endpoint
	UpdateUserMetadata endpoint.Endpoint
	GetUsers endpoint.Endpoint
	GetUsersPaged endpoint.Endpoint
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

	var authenticateOAuthRedirectEndpoint endpoint.Endpoint
	{
		authenticateOAuthRedirectEndpoint = AuthenticateOAuthRedirectEndpoint(svc)

		authenticateOAuthRedirectEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "authenticateOauthRedirect"))(authenticateOAuthRedirectEndpoint)
	}

	var authenticateAccessTokenEndpoint endpoint.Endpoint
	{
		authenticateAccessTokenEndpoint = AuthenticateAccessTokenEndpoint(svc)

		authenticateAccessTokenEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "authenticateGoogleAccessToken"))(authenticateAccessTokenEndpoint)
	}

	var updateUserEndpoint endpoint.Endpoint
	{
		updateUserEndpoint = UpdateUserEndpoint(svc)

		updateUserEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "updateUser"))(updateUserEndpoint)
	}

	var updateUserMetadataEndpoint endpoint.Endpoint
	{
		updateUserMetadataEndpoint = UpdateUserMetadataEndpoint(svc)

		updateUserMetadataEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "updatedUserMetadata"))(updateUserMetadataEndpoint)
	}

	var getUsersEndpoint endpoint.Endpoint
	{
		getUsersEndpoint = GetUsersEndpoint(svc)

		getUsersEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "getUsersEndpoint"))(getUsersEndpoint)
	}

	var getUsersPagedEndpoint endpoint.Endpoint
	{
		getUsersPagedEndpoint = GetUsersPagedEndpoint(svc)
		getUsersPagedEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "getUsersPagedEndpoint"))(getUsersPagedEndpoint)
	}

	return Set{
		ValidateLogin: validateLoginEndpoint,
		SignupUser: signupUserEndpoint,
		UpdateUserRole: updateUserRoleEndpoint,
		GenerateAPIKey: generateAPIKeyEndpoint,
		LoginWithAPIKey: loginWithAPIKeyEndpoint,
		RegisterOAuth: registerOAuthEndpoint,
		AuthenticateOAuthRedirect: authenticateOAuthRedirectEndpoint,
		AuthenticateAccessToken: authenticateAccessTokenEndpoint,
		UpdateUser: updateUserEndpoint,
		UpdateUserMetadata: updateUserMetadataEndpoint,
		GetUsers: getUsersEndpoint,
		GetUsersPaged: getUsersPagedEndpoint,
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
			"scopes": req.Scopes,
		}

		metadataBytes, jsonErr := json.Marshal(metadata)

		if jsonErr != nil {
			return nil, jsonErr
		}

		newOAuthProvider := usermodels.OAuthProvider{
			Name: req.ProviderName,
			Metadata: metadataBytes,
			TenantIdentifier: req.TenantIdentifier,
		}

		url, err := s.RegisterOAuth(ctx, newOAuthProvider, req.Scopes, req.Host)

		return RegisterGoogleOAuthResponse{url, err}, err

	}
}

func AuthenticateOAuthRedirectEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthenticateGoogleOAuthRedirectRequest)

		if req.State != os.Getenv("STATE") {
			return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
		}

		jwt, err := s.AuthenticateOAuthRedirect(ctx, req.Code, req.ProviderName, req.TenantIdentifier, req.Host)

		return AuthenticateGoogleOAuthRedirectResponse{jwt, err}, err
	}

}

func AuthenticateAccessTokenEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthenticateAccessTokenRequest)

		jwt, user, err := s.AuthenticateAccessToken(ctx, req.AccessToken, req.TenantIdentifier, req.ProviderName)

		return AuthenticateAccessTokenResponse{jwt, user, err}, err
	}
}

func UpdateUserEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateUserRequest)

		user, err := s.UpdateUser(ctx, req.NewUser)

		return UpdateUserResponse{user, err}, err
	}
}

func UpdateUserMetadataEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateUserMetadataRequest)

		user, err := s.UpdateUserMetadata(ctx, req.UpdatedUser)

		return UpdateUserMetadataResponse{user, err}, err
	}
}

func GetUsersEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUsersRequest)

		users, err := s.GetUsers(ctx, req.UserSearch)

		return GetUsersResponse{users, err}, err
	}
}

func GetUsersPagedEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(commonmodels.PagedRequest[GetUsersRequest])

		pagedUsers, err := s.GetUsersPaged(ctx, req.Entity.UserSearch, req.Page, req.Size)

		return pagedUsers, err

	}
}

