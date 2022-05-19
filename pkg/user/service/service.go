package userservice

import "context"

type UserService interface {
	GenerateToken(ctx context.Context) string
}

type userService struct {
} //class-like skeleton in Go

func NewUserService() UserService { //makes userService struct implement UserService interface
	return &userService{} //returns an address which points to userService to make changes in original memory address
}

func (s *userService) GenerateToken(_ context.Context) string {
	return "dummyJWT"
}

func init(){}
