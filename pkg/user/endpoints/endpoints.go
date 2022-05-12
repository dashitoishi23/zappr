package userendpoint

import (
	"context"
	"errors"
	"time"

	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
)

type Set struct {
	GenerateToken endpoint.Endpoint
} //defines all endpoints as having type Endpoint, provided by go-kit

func New(svc userservice.UserService) Set {
	var generateTokenEndpoint endpoint.Endpoint
	{
		generateTokenEndpoint = GenerateTokenEndpoint(svc)
		generateTokenEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(generateTokenEndpoint)
	}

	return Set{
		GenerateToken: generateTokenEndpoint,
	}
}

func GenerateTokenEndpoint(s userservice.UserService) endpoint.Endpoint{
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