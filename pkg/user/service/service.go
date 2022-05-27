package userservice

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserService interface {
	GenerateJWTToken(ctx context.Context, userEmail string) (string, error)
	ValidateLogin(ctx context.Context, jwt string) (bool)
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

func (s *userService) ValidateLogin(_ context.Context, jwtToken string) bool {
	jwtToken = strings.Split(jwtToken, " ")[1]

	parsedToken, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC) 
		
		if !ok {
			return nil, fmt.Errorf("unauthorized attempt")
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})

	if err != nil {
		return false
	}

	return parsedToken.Valid

}
