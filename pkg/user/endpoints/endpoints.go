package userendpoint

import (
	"context"
	"fmt"

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
	return func(ctx context.Context, request interface{}) (response interface{}, err error){
		_, ok:= request.(GenerateTokenRequest)
		s := s.GenerateToken(ctx)
		if !ok {
			fmt.Println(err.Error())
			return GenerateTokenResponse{s, err.Error()}, nil
		}
		return GenerateTokenResponse{s, ""}, nil
	}
}

type GenerateTokenRequest struct {

} //strongly typed request object

type GenerateTokenResponse struct {
	s string `json:"s"`
	err string `json:"err,omitempty"`
} //strongly typed response object

func init(){}