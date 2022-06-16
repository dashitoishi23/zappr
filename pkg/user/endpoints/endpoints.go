package userendpoint

import (
	"context"

	util "dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/util"
	models "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	GenerateToken endpoint.Endpoint
	ValidateLogin endpoint.Endpoint
	SignupUser endpoint.Endpoint
} //defines all endpoints as having type Endpoint, provided by go-kit

func New(svc userservice.UserService, logger log.Logger) Set {

	var generateTokenEndpoint endpoint.Endpoint
	{
		generateTokenEndpoint = GenerateTokenEndpoint(svc)

		generateTokenEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "generateToken"))(generateTokenEndpoint)

	}

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

	return Set{
		GenerateToken: generateTokenEndpoint,
		ValidateLogin: validateLoginEndpoint,
		SignupUser: signupUserEndpoint,
	}
}

func GenerateTokenEndpoint(s userservice.UserService) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (interface{}, error){
		req := request.(GenerateTokenRequest)
		s, err := s.GenerateJWTToken(ctx, req.UserEmail)
		
		if err !=nil {
			return GenerateTokenResponse{"", err.Error()}, err
		}
		return GenerateTokenResponse{s, ""}, nil
	}
}

func ValidateLoginEndpoint(s userservice.UserService) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ValidateLoginRequest)
		resp := s.ValidateLogin(ctx, req.Token)

		return ValidateLoginResponse{resp}, nil
	}
}

func SignupUserEndpoint(s userservice.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SignupUserRequest)
		req.NewUser.InitFields()
		resp, err := s.SignupUser(ctx, req.NewUser)

		return SignupUserResponse{resp}, err
	}
}

type GenerateTokenRequest struct {
	UserEmail string `json:"useremail"`
} //strongly typed request object

type GenerateTokenResponse struct {
	JWT string `json:"jwt"`
	Err string `json:"err,omitempty"`
} //strongly typed response object

type ValidateLoginRequest struct {
	Token string `json:"token"`
}

type ValidateLoginResponse struct {
	IsValid bool `json:"isvalid"`
}

type SignupUserRequest struct {
	NewUser models.User `json:"newuser"`
}

type SignupUserResponse struct {
	NewUser models.User `json:"newuser"`
}