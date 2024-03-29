package userendpoint

import (
	models "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	userrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/userrole/models"
	"github.com/lib/pq"
)

type GenerateTokenRequest struct {
	UserEmail string `json:"useremail"`
} //strongly typed request object

type GenerateTokenResponse struct {
	JWT string `json:"jwt"`
	Err error  `json:"-"`
} //strongly typed response object

func (g GenerateTokenResponse) Failed() error { return g.Err }

type ValidateLoginRequest struct {
	CurrentUser models.UserLogin `jwt:"currentuser"`
}

type ValidateLoginResponse struct {
	Jwt string  `json:"jwt"`
	Err     error `json:"-"`
}

//Implements the endpoint.Failer interface
func (v ValidateLoginResponse) Failed() error { return v.Err }

type SignupUserRequest struct {
	NewUser models.User `json:"newuser"`
}

type SignupUserResponse struct {
	NewUser models.User `json:"newuser"`
	Err     error       `json:"-"`
}

func (s SignupUserResponse) Failed() error { return s.Err }

type UpdateUserRoleRequest struct {
	MasterRoleIdentifier string `json:"masterRoleIdentifier"`
	UserIdentifier string `json:"userIdentifier"`
}

type UpdateUserRoleResponse struct {
	UpdatedUserRole userrolemodels.UserRole `json:"updatedUserRole"`
	Err 	error 		`json:"-"`
}

func (u UpdateUserRoleResponse) Failed() error { return u.Err }

type GenerateAPIKeyRequest struct {
	APIKeyName string `json:"apiKeyName"`
}

type GenerateAPIKeyResponse struct {
	APIKey string `json:"apiKey"`
	Err 	error `json:"-"`
}

func (g GenerateAPIKeyResponse) Failed() error { return g.Err }

type LoginWithAPIKeyRequest struct {
	APIKey string `json:"apiKey"`
}

type LoginWithAPIKeyResponse struct {
	JWT string `json:"jwt"`
	Err error `json:"-"`
}

func (l LoginWithAPIKeyResponse) Failed() error { return l.Err }

type RegisterOAuthRequest struct {
	ProviderName string `json:"providerName"`
	ClientID string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURI string `json:"redirectURI"`
	Scopes pq.StringArray `json:"scopes"`
	Host string `json:"host"`
	TenantIdentifier string `json:"tenantIdentifier"`
}

type RegisterGoogleOAuthResponse struct {
	AuthDialogURL string `json:"authDialogURL"`
	Err error `json:"-"`
}

func (r RegisterGoogleOAuthResponse) Failed() error { return r.Err }

type AuthenticateGoogleOAuthRedirectRequest struct {
	Code string `json:"code"`
	State string `json:"state"`
	ProviderName string `json:"provideName"`
	Host string `json:"host"`
	TenantIdentifier string `json:"tenantIdentifier"`
}

type AuthenticateGoogleOAuthRedirectResponse struct {
	Jwt string `json:"jwt"`
	Err error `json:"error"`
}

func (a *AuthenticateGoogleOAuthRedirectResponse) Failed() error { return a.Err }

type AuthenticateAccessTokenRequest struct {
	AccessToken string `json:"accessToken"`
	TenantIdentifier string `json:"tenantIdentifier"`
	ProviderName string `json:"providerName"`
}

type AuthenticateAccessTokenResponse struct {
	Jwt string `json:"jwt"`
	User models.User `json:"user"`
	Err error `json:"-"`
}

func (a *AuthenticateAccessTokenResponse) Failed() error { return a.Err }

type UpdateUserRequest struct {
	NewUser models.UpdateUser `json:"newUser"`
}

type UpdateUserResponse struct {
	NewUser models.User `json:"newUser"`
	Err error `json:"-"`
}

func (u *UpdateUserResponse) Failed() error { return u.Err }

type UpdateUserMetadataRequest struct {
	UpdatedUser models.UpdateUserMetadata `json:"updatedUser"`
}

type UpdateUserMetadataResponse struct {
	UpdatedUser models.User `json:"updatedUser"`
	Err error `json:"-"`
}

func (u *UpdateUserMetadataResponse) Failed() error { return u.Err }

type GetUsersRequest struct {
	UserSearch models.SearchableUser `json:"userSearch"`
}

type GetUsersResponse struct {
	Users []models.User `json:"users"`
	Err error `json:"-"`
}

func (g *GetUsersResponse) Failed() error { return g.Err }

type GetUsersByMetadataRequest struct {
	Query map[string]interface{} `json:"query"`
}

type GetUsersByMetadataResponse struct {
	Users []models.User `json:"users"`
	Err error `json:"-"`
}

func (g *GetUsersByMetadataResponse) Failed() error { return g.Err }

type UpdateCurrentUserRequest struct {
	NewUser models.UpdateCurrentUser `json:"newUser"`
}

type UpdateCurrentUserResponse struct {
	UpdatedUser models.User `json:"updatedUser"`
	Err error `json:"-"`
}

func (u *UpdateCurrentUserResponse) Failed() error { return u.Err }

type GetCurrentUserDetailsRequest struct {

}

type GetCurrentUserDetailsResponse struct {
	User models.User `json:"user"`
	Err error `json:"-"`
}

func (g *GetCurrentUserDetailsResponse) Failed() error {return g.Err}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordResponse struct {
	IsPasswordChanged bool `json:"isPasswordChanged"`
	Err error `json:"-"`
}

func (c *ChangePasswordResponse) Failed() error { return c.Err }
