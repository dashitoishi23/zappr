package tenanttransports

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	tenantendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/endpoints"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHandler(endpoints tenantendpoint.Set) []commonmodels.HttpServerConfig {
	var tenantServers []commonmodels.HttpServerConfig

	createHandler := httptransport.NewServer(
		endpoints.CreateTenant,
		util.DecodeHTTPGenericRequest[tenantendpoint.CreateTenantRequest],
		util.EncodeHTTPGenericResponse,
	)

	findFirstHandler := httptransport.NewServer(
		endpoints.FindFirstTenant,
		util.DecodeHTTPGenericRequest[tenantendpoint.FindFirstTenantRequest],
		util.EncodeHTTPGenericResponse,
	)

	getAllTenantsHandler := httptransport.NewServer(
		endpoints.GetAllTenants,
		util.DecodeHTTPGenericRequest[tenantendpoint.GetAllTenantsRequest],
		util.EncodeHTTPGenericResponse,
	)

	findTenantsHandler := httptransport.NewServer(
		endpoints.FindTenants,
		util.DecodeHTTPGenericRequest[tenantendpoint.FindTenantsRequest],
		util.EncodeHTTPGenericResponse,
	)

	return append(tenantServers, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: createHandler,
		Route: "/tenant",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: findFirstHandler,
		Route: "/tenant",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: getAllTenantsHandler,
		Route: "/tenant/all",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: findTenantsHandler,
		Route: "/tenants/find",
		Methods: []string{"GET"},
	})
}