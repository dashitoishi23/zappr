package userservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	constants "dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	masterrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/role/models"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	models "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	userrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/userrole/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/golang-jwt/jwt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/validator.v2"
)

type UserService interface {
	SignupUser(ctx context.Context, newUser models.User) (models.User, error)
	LoginUser(ctx context.Context, currentUser models.UserLogin) (string, error)
	UpdateUserRole(ctx context.Context, roleIdentifier string, 
		userIdentifier string) (userrolemodels.UserRole, error)
	FindUserById(ctx context.Context, identifier string) (models.User, error)
	GenerateAPIKey(ctx context.Context, apiKeyName string) (string, error)
	LoginWithAPIKey (ctx context.Context, apiKey string) (string, error)
	RegisterGoogleOAuth(ctx context.Context, newOAuthProvider models.OAuthProvider) (string, error)
}

type userService struct {
	repository repository.IRepository[models.User]
	roleRepository repository.IRepository[masterrolemodels.Role]
	userRoleRepository repository.IRepository[userrolemodels.UserRole]
	tenantRepository repository.IRepository[tenantmodels.Tenant]
	apiKeyRepository repository.IRepository[models.APIKey]
	oAuthProviderRepository repository.IRepository[models.OAuthProvider]
} //class-like skeleton in Go

func NewUserService(repository repository.IRepository[models.User], 
	roleRepository repository.IRepository[masterrolemodels.Role], 
	userRoleRepository repository.IRepository[userrolemodels.UserRole], 
	tenantRepository repository.IRepository[tenantmodels.Tenant], 
	apiKeyRepository repository.IRepository[models.APIKey], 
	oAuthProviderRepository repository.IRepository[models.OAuthProvider]) UserService { //makes userService struct implement UserService interface
	return &userService{
		repository: repository,
		roleRepository: roleRepository,
		userRoleRepository: userRoleRepository,
		tenantRepository: tenantRepository,
		apiKeyRepository: apiKeyRepository,
		oAuthProviderRepository: oAuthProviderRepository,

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

	tenant, err := s.tenantRepository.FindFirst(&tenantmodels.SearchableTenant{
		Identifier: newUser.TenantIdentifier,
	})

	if err != nil {
		return newUser, err
	}

	if tenant.AdminEmail == newUser.Email {
	role, roleErr := s.roleRepository.FindFirst(masterrolemodels.SearchableRole{
	Name: "Admin",
	TenantIdentifier: newUser.TenantIdentifier,
	})

	if roleErr != nil {
		return newUser, roleErr
	}

	var scopes []string

	role.Scopes.Scan(pq.Array(&scopes))
	
	newUser.Role = userrolemodels.UserRole{
		UserIdentifier: newUser.Identifier,
		RoleIdentifier: role.Identifier,
		Scopes: role.Scopes,
	}

	newUser.Role.InitFields()

	tx := s.repository.Add(newUser)


	if tx.Error != nil {
		return newUser, tx.Error
	}

	return newUser, nil

	}

	role, roleErr := s.roleRepository.FindFirst(masterrolemodels.SearchableRole{
	Name: "Normal User",
	TenantIdentifier: newUser.TenantIdentifier,
	})

	if roleErr != nil {
		return newUser, roleErr
	}

	var scopes []string

	role.Scopes.Scan(pq.Array(&scopes))
	
	newUser.Role = userrolemodels.UserRole{
		UserIdentifier: newUser.Identifier,
		RoleIdentifier: role.Identifier,
		Scopes: role.Scopes,
	}

	newUser.Role.InitFields()

	tx := s.repository.Add(newUser)


	if tx.Error != nil {
		return newUser, tx.Error
	}

	return newUser, nil

}

func (s *userService) LoginUser (ctx context.Context, currentUser models.UserLogin) (string, error) {
	if errs:= validator.Validate(currentUser); errs != nil {
		return "", errors.New(constants.INVALID_MODEL)
	}

	existingUser, err := s.repository.FindFirstByAssociation("Role", &models.SearchableUser{
		Email: currentUser.Email,
		TenantIdentifier: currentUser.TenantIdentifier,
	})

	if err != nil {
		if err.Error() == constants.RECORD_NOT_FOUND {
			return "", errors.New(constants.RECORD_NOT_FOUND)
		}
		return "", err
	}

	fmt.Print(existingUser.Role.Scopes.Value())

	hashErr := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(currentUser.Password))

	if hashErr == nil {
		jwt, _ := s.generateJWTToken(ctx, existingUser.Email, existingUser.TenantIdentifier, existingUser.Identifier, 
		existingUser.Role.Scopes)

		conf := &oauth2.Config{
			ClientID:     "548368247582-nrt0g3llbjfgf7q8lakije8kh713rla1.apps.googleusercontent.com",
			ClientSecret: "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",
			RedirectURL:  "http://localhost:8000/oauth/google",
			Scopes: []string{
				"https://www.googleapis.com/auth/bigquery",
				"https://www.googleapis.com/auth/blogger",
			},
			Endpoint: google.Endpoint,
		}
		url := conf.AuthCodeURL("state")
		fmt.Printf("Visit the URL for the auth dialog: %v", url)

		tok, err := conf.Exchange(ctx, "4/0AdQt8qgFiScWMJihi4QITPY96yXK8EyZifucM4uTvNi6oXBSifUDVBTqL_jkIMqZYmfP-A")
		if err != nil {
			log.Fatal(err)
		}
		client := conf.Client(ctx, tok)

		fmt.Print(client)

		return jwt, nil
	}

	return "", errors.New(constants.UNAUTHORIZED_ATTEMPT)


}

func (s *userService) UpdateUserRole(ctx context.Context, roleIdentifier string, 
	userIdentifier string) (userrolemodels.UserRole, error) {
	
	existingUser, err := s.FindUserById(ctx, userIdentifier)

	existingRole, roleErr := s.roleRepository.FindFirst(&masterrolemodels.SearchableRole{
		Identifier: roleIdentifier,
	})

	if err != nil {
		return existingUser.Role, err
	}

	if roleErr != nil {
		return existingUser.Role, roleErr
	}

	updatedUserRole, roleErr := s.userRoleRepository.Update(userrolemodels.UserRole{
		Identifier: existingUser.Role.Identifier,
		RoleIdentifier: roleIdentifier,
		UserIdentifier: userIdentifier,
		Scopes: existingRole.Scopes,
		CreatedOn: existingUser.Role.CreatedOn,
		ModifiedOn: time.Now(),
	})

	return updatedUserRole, roleErr	
}

func (s *userService) FindUserById(ctx context.Context, identifier string) (models.User, error) {
	res, err := s.repository.FindFirstByAssociation("Role", &models.SearchableUser{
		Identifier: identifier,
		TenantIdentifier: ctx.Value("requestScope").(commonmodels.RequestScope).UserTenant,
	})

	return res, err
}

func (s *userService) GenerateAPIKey(ctx context.Context, apiKeyName string) (string, error) {
	requestScope := ctx.Value("requestScope").(commonmodels.RequestScope)

	newAPIKey := models.APIKey{
		Secret: util.RandStringRunes(),
		Scopes: requestScope.UserScopes,
		UserIdentifier: requestScope.UserIdentifier,
		Name: apiKeyName,
	}

	newAPIKey.InitFields(requestScope)

	tx := s.apiKeyRepository.Add(newAPIKey)

	if tx.Error != nil {
		return "", tx.Error
	}

	return newAPIKey.Secret, nil

}

func (s *userService) LoginWithAPIKey(ctx context.Context, apiKey string) (string, error) {
	requestScope := ctx.Value("requestScope").(commonmodels.RequestScope)

	_, err:= s.apiKeyRepository.FindFirst(models.APIKey{
		Secret: apiKey,
		TenantIdentifier: requestScope.UserTenant,
		UserIdentifier: requestScope.UserIdentifier,
	})

	if err != nil {
		return "", err
	}

	existingUser, err := s.repository.FindFirstByAssociation("Role", models.SearchableUser{
		Identifier: requestScope.UserIdentifier,
		TenantIdentifier: requestScope.UserTenant,
	})

	if err != nil {
		return "", err
	}

	return s.generateJWTToken(ctx, existingUser.Email, existingUser.TenantIdentifier, existingUser.Identifier, 
	existingUser.Role.Scopes)
}

func (s *userService) generateJWTToken(_ context.Context, userEmail string, tenantIdentifier string, 
	userIdentifier string, userScopes pq.StringArray) (string, error) {
	if userEmail == "" {
		return "", errors.New(constants.INVALID_MODEL)
	}

	claims := commonmodels.JWTClaims{
		userEmail,
		tenantIdentifier,
		userIdentifier,
		userScopes,
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

func (s *userService) RegisterGoogleOAuth(ctx context.Context, newOAuthProvider models.OAuthProvider) (string, error) {
	tx := s.oAuthProviderRepository.Add(newOAuthProvider)

	if tx.Error != nil {
		return "", tx.Error
	}

	conf := &oauth2.Config{
		ClientID:     newOAuthProvider.Metadata["client_id"].(string),
		ClientSecret: newOAuthProvider.Metadata["client_secret"].(string),
		RedirectURL:  newOAuthProvider.Metadata["redirect_uri"].(string),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	stateToken := os.Getenv("GOOGLE_OAUTH_STATE")
	
	url := conf.AuthCodeURL(stateToken)

	return url, nil

}



