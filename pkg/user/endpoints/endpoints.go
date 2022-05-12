package userendpoint

import (
	"context"
	"errors"

	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	GenerateToken endpoint.Endpoint
} //defines all endpoints as having type Endpoint, provided by go-kit

func GenerateToken(s userservice.UserService) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error){
		_, ok:= request.(GenerateTokenRequest)
		s := s.GenerateToken(ctx)
		if !ok {
			err:= errors.New("Invalid request")
		}
		return GenerateTokenResponse{s: s, err: err}, nil
	}
}

type GenerateTokenRequest struct {

} //strongly typed request object

type GenerateTokenResponse struct {
	s string `json:"s"`
	err error `json:"err"`
} //strongly typed response object