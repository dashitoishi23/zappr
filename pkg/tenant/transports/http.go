package tenanttransports

import (
	"context"
	"encoding/json"
	"net/http"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/util"
	tenantendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/endpoints"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHandler(endpoints tenantendpoint.Set) []commonmodels.HttpServerConfig {
	var tenantServers []commonmodels.HttpServerConfig

	createHandler := httptransport.NewServer(
		endpoints.CreateTenant,
		decodeCreateTenantRequest,
		util.EncodeHTTPGenericResponse,
	)

	return append(tenantServers, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: createHandler,
		Route: "/tenant",
		Methods: []string{"POST"},
	})
}

func decodeCreateTenantRequest(_ context.Context, r *http.Request) (interface{}, error){
	var req tenantendpoint.CreateTenantRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return req, err
	}

	return req, nil

}