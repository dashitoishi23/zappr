package userservice

import (
	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	GenerateToken endpoint.Endpoint
} //defines all endpoints as having type Endpoint, provided by go-kit

func GenerateToken(s userservice.UserService)