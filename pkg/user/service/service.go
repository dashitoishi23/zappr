package userservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	constants "dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	masterrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/role/models"
	models "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	userrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/userrole/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

type UserService interface {
	SignupUser(ctx context.Context, newUser models.User) (models.User, error)
	LoginUser(ctx context.Context, currentUser models.UserLogin) (string, error)
}

type userService struct {
	repository repository.IRepository[models.User]
	roleRepository repository.IRepository[masterrolemodels.Role]
} //class-like skeleton in Go

func NewUserService(repository repository.IRepository[models.User], 
	roleRepository repository.IRepository[masterrolemodels.Role]) UserService { //makes userService struct implement UserService interface
	return &userService{
		repository: repository,
		roleRepository: roleRepository,
	} //returns an address which points to userService to make changes in original memory address
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

	role, roleErr := s.roleRepository.FindFirst(masterrolemodels.SearchableRole{
	Name: "Normal User",
	TenantIdentifier: newUser.TenantIdentifier,
	})

	if roleErr != nil {
		return newUser, roleErr
	}
	
	newUser.Role = userrolemodels.UserRole{
		UserIdentifier: newUser.Identifier,
		RoleIdentifier: role.Identifier,
	}

	newUser.Role.InitFields()

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
		TenantIdentifier: currentUser.TenantIdentifier,
	})

	if err != nil {
		if err.Error() == constants.RECORD_NOT_FOUND {
			return "", errors.New(constants.RECORD_NOT_FOUND)
		}
		return "", errors.New(constants.UNAUTHORIZED_ATTEMPT)
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(currentUser.Password))

	if hashErr == nil {
		jwt, _ := s.generateJWTToken(ctx, existingUser.Email, existingUser.TenantIdentifier, existingUser.Identifier)

		return jwt, nil
	}

	return "", errors.New(constants.UNAUTHORIZED_ATTEMPT)


}

func (s *userService) generateJWTToken(_ context.Context, userEmail string, tenantIdentifier string, 
	userIdentifier string) (string, error) {
	if userEmail == "" {
		return "", errors.New(constants.INVALID_MODEL)
	}

	claims := commonmodels.JWTClaims{
		userEmail,
		tenantIdentifier,
		userIdentifier,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer: "zappr",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	signingKey := []byte(os.Getenv("JWT_SIGNING_KEY"))

	ss, err := token.SignedString(signingKey)

	return ss, err
}