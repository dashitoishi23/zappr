package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	constants "dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
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

	if !tokenValidator(authtoken) {
		response["isValid"] = false

		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(response)

		return
	}

	next(w, r)

}

func tokenValidator(jwtToken string) bool {
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
