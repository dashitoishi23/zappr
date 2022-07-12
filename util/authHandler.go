package util

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	constants "dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"github.com/golang-jwt/jwt"
)

var AuthHandler = func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authtoken := r.Header.Get("Authorization")
	
	response := make(map[string]bool)

	if authtoken == "" {

		response["isValid"] = false

		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(response)

		return

	}

	resp, ctx := tokenValidator(authtoken)

	if !resp {
		response["isValid"] = false

		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(response)

		return
	}

	next(w, r.WithContext(ctx))

}

func tokenValidator(jwtToken string) (bool, context.Context) {
	ctx := context.Background()

	if jwtToken == "" {
		return false, ctx
	}

	jwtToken = strings.Split(jwtToken, " ")[1]

	parsedToken, err := jwt.ParseWithClaims(jwtToken, &commonmodels.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC) 
		
		if !ok {
			return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
		}

	//Checks against signing algorithm forgery which involves switching to RSA and confusing the verification
	if t.Method.Alg() != jwt.SigningMethodHS256.Name {
		return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
	} 

		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})

	if err != nil {
		return false, ctx
	}

	if claims, ok := parsedToken.Claims.(*commonmodels.JWTClaims); ok && parsedToken.Valid {
		ctx = context.WithValue(ctx, "requestScope", commonmodels.RequestScope{
			UserTenant: claims.UserTenant,
			UserIdentifier: claims.UserIdentifier,
			UserScopes: claims.UserScopes,
		})

		return parsedToken.Valid, ctx
	}

	return false, ctx
}
