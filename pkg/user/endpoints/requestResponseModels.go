package userendpoint

import models "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"

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