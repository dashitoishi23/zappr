package usermodels

import "time"

type SearchableUser struct {
	Identifier       string     `json:"identifier"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Password         string     `json:"password"`
	IsExternalOAuthUser         bool       `json:"isExternalOAuthUser"`
	OAuthProvider    string 	`json:"oAuthProvider"`
	Locale           string     `json:"locale"`
	TenantIdentifier string     `json:"tenantIdentifier"`
	CreatedOn        time.Time  `json:"createdOn"`
	ModifiedOn       *time.Time `json:"modifiedOn"`
}