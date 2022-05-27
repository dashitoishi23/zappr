package userservice

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserService interface {
	GenerateJWTToken(ctx context.Context, userEmail string) (string, error)
}

type userService struct {
} //class-like skeleton in Go

func NewUserService() UserService { //makes userService struct implement UserService interface
	return &userService{} //returns an address which points to userService to make changes in original memory address
}

type zapprJWTClaims struct {
	UserEmail string `json: "userEmail"`
	jwt.StandardClaims
}

func (s *userService) GenerateJWTToken(_ context.Context, userEmail string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, zapprJWTClaims{
		userEmail,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer: "zappr",
		},
	})
	
	signingKey := []byte(os.Getenv("JWT_SIGNING_KEY"))

	ss, err := token.SignedString(signingKey)

	return ss, err
}
