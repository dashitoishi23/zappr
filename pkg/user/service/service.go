package userservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	constants "dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	models "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	state "dev.azure.com/technovert-vso/Zappr/_git/Zappr/state"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

type UserService interface {
	generateJWTToken(ctx context.Context, userEmail string) (string, error)
	SignupUser(ctx context.Context, newUser models.User) (models.User, error)
	LoginUser(ctx context.Context, currentUser models.UserLogin) (string, error)
}

type userService struct {
	repository repository.IRepository[models.User]
} //class-like skeleton in Go

func NewUserService(repository repository.IRepository[models.User]) UserService { //makes userService struct implement UserService interface
	return &userService{
		repository: repository,
	} //returns an address which points to userService to make changes in original memory address
}

type zapprJWTClaims struct {
	UserEmail string `json:"userEmail"`
	jwt.StandardClaims
}

func (s *userService) SignupUser(_ context.Context, newUser models.User) (models.User, error) {
	if errs := validator.Validate(newUser); errs != nil {
		return newUser, errors.New(constants.INVALID_MODEL)
	}

	if !newUser.IsADUser {
		userPwd := []byte(newUser.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(userPwd, 10)

		if err != nil {
			return newUser, err
		}

		newUser.Password = string(hashedPassword)
	}

	tx := s.repository.Add(newUser)

	if tx.Error != nil {
		return newUser, tx.Error
	}

	fmt.Print(tx.RowsAffected)

	return newUser, nil

}

func (s *userService) LoginUser (ctx context.Context, currentUser models.UserLogin) (string, error) {
	if errs:= validator.Validate(currentUser); errs != nil {
		return "", errors.New(constants.INVALID_MODEL)
	}

	existingUser, err := s.repository.FindFirst(&models.SearchableUser{
		Email: currentUser.Email,
	})

	if err != nil {
		if err.Error() == constants.RECORD_NOT_FOUND {
			return "", errors.New(constants.RECORD_NOT_FOUND)
		}
		return "", errors.New(constants.UNAUTHORIZED_ATTEMPT)
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(currentUser.Password))

	if hashErr == nil {
		jwt, _ := s.generateJWTToken(ctx, existingUser.Email)
		state := state.GetState()

		state.SetUserContext(commonmodels.UserContext{
			UserTenant: existingUser.TenantIdentifier,
			UserIdentifier: existingUser.Identifier,
		})

		return jwt, nil
	}

	return "", errors.New(constants.UNAUTHORIZED_ATTEMPT)


}

func (s *userService) generateJWTToken(_ context.Context, userEmail string) (string, error) {
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