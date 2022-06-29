package commonmodels

import "github.com/golang-jwt/jwt"

type JWTClaims struct {
	UserEmail string `json:"userEmail"`
	UserTenant string `json:"userTenant"`
	UserIdentifier string `json:"userIdentifier"`
	jwt.StandardClaims
}