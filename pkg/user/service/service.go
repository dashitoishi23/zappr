package userservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	constants "dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/constants"
	util "dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/util"
	models "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type UserService interface {
	GenerateJWTToken(ctx context.Context, userEmail string) (string, error)
	ValidateLogin(ctx context.Context, jwt string) (bool)
	SignupUser(ctx context.Context, newUser models.User) (models.User, error)
}

type userService struct {
	db *gorm.DB
	repository util.IRepository[models.User]
} //class-like skeleton in Go

func NewUserService(database *gorm.DB) UserService { //makes userService struct implement UserService interface
	return &userService{
		db : database,
		repository: util.Repository[models.User](database),
	} //returns an address which points to userService to make changes in original memory address
}

type zapprJWTClaims struct {
	UserEmail string `json:"userEmail"`
	jwt.StandardClaims
}

func (s *userService) GenerateJWTToken(_ context.Context, userEmail string) (string, error) {
	if userEmail == "" {
		return "", errors.New(constants.INVALID_MODEL)
	}

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
	if jwtToken == "" {
		return false
	}

	jwtToken = strings.Split(jwtToken, " ")[1]

	parsedToken, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC) 
		
		if !ok {
			return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})

	if err != nil {
		return false
	}

	return parsedToken.Valid

}

func (s *userService) SignupUser(_ context.Context, newUser models.User) (models.User, error) {
	if errs := validator.Validate(newUser); errs != nil {
		return newUser, errors.New(constants.INVALID_MODEL)
	}

	newUser.Identifier = uuid.New().String()

	tx := s.repository.Add(newUser)

	if tx.Error != nil {
		return newUser, tx.Error
	}

	fmt.Print(tx.RowsAffected)

	return newUser, nil
}
