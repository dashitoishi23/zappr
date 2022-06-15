package tenanttransports

import (
	"context"
	"encoding/json"
	"net/http"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/util"
	tenantendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/endpoints"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHandler(endpoints tenantendpoint.Set) http.Handler {
	createHandler := httptransport.NewServer(
		endpoints.CreateTenant,
		decodeCreateTenantRequest,
		util.EncodeHTTPGenericResponse,
	)

	r := mux.NewRouter()

	r.Handle("/tenant", createHandler).Methods("POST")

	return r
}

func decodeCreateTenantRequest(_ context.Context, r *http.Request) (interface{}, error){
	var req tenantendpoint.CreateTenantRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return req, err
	}

	return req, nil

}