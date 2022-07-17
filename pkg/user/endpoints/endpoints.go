package userendpoint

import (
	"context"

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

	return Set{
		ValidateLogin: validateLoginEndpoint,
		SignupUser: signupUserEndpoint,
		UpdateUserRole: updateUserRoleEndpoint,
		GenerateAPIKey: generateAPIKeyEndpoint,
		LoginWithAPIKey: loginWithAPIKeyEndpoint,
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