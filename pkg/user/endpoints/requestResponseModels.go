package userendpoint

import (
	models "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	userrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/userrole/models"
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

type RegisterGoogleOAuthRequest struct {
	ClientID string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURI string `json:"redirectURI"`
}

type RegisterGoogleOAuthResponse struct {
	AuthDialogURL string `json:"authDialogURL"`
	Err error `json:"-"`
}

func (r RegisterGoogleOAuthResponse) Failed() error { return r.Err }

type AuthenticateGoogleOAuthRedirectRequest struct {
	Code string `json:"code"`
	State string `json:"state"`
}

type AuthenticateGoogleOAuthRedirectResponse struct {
	Jwt string `json:"jwt"`
	Err error `json:"error"`
}

func (a *AuthenticateGoogleOAuthRedirectResponse) Failed() error { return a.Err }

