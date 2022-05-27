package userendpoint

import (
	"context"

	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	GenerateToken endpoint.Endpoint
} //defines all endpoints as having type Endpoint, provided by go-kit

func New(svc userservice.UserService) Set {
	return Set{
		GenerateToken: GenerateTokenEndpoint(svc),
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

func (s *Set) GenerateJWTToken(ctx context.Context) string {
	resp, _ := s.GenerateToken(ctx, GenerateTokenRequest{})

	getTokenResp := resp.(GenerateTokenResponse)

	return getTokenResp.JWT

}

type GenerateTokenRequest struct {
	UserEmail string `json: "useremail"`
} //strongly typed request object

type GenerateTokenResponse struct {
	JWT string `json:"jwt"`
	Err string `json:"err,omitempty"`
} //strongly typed response object