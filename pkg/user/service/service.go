package userservice

import (
	"context"
	"fmt"
)

type UserService interface {
	GenerateJWTToken(ctx context.Context) string
}

type userService struct {
} //class-like skeleton in Go

func NewUserService() UserService { //makes userService struct implement UserService interface
	return &userService{} //returns an address which points to userService to make changes in original memory address
}

func (s *userService) GenerateJWTToken(_ context.Context) string {
	fmt.Println("Services")
	return "dummyJWT"
}
