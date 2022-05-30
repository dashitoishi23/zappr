package util

import (
	"context"
	"encoding/json"
	"net/http"

	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
)

var AuthHandler = func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authtoken := r.Header.Get("Authorization")

	ctx := context.Background()
	
	response := make(map[string]bool)

	if authtoken == "" {

		response["isValid"] = false

		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(response)

		return

	}

	if !userservice.NewUserService().ValidateLogin(ctx, authtoken) {
		response["isValid"] = false

		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(response)

		return
	}

	next(w, r)

}
