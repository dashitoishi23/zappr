package tenanttransports

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	tenantendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/endpoints"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func NewHandler(endpoints tenantendpoint.Set, logger log.Logger) []commonmodels.HttpServerConfig {
	var tenantServers []commonmodels.HttpServerConfig

	serverOptions := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(util.ErrorEncoder),
	}

	createHandler := httptransport.NewServer(
		endpoints.CreateTenant,
		util.DecodeHTTPGenericRequest[tenantendpoint.CreateTenantRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...,
	)

	findFirstHandler := httptransport.NewServer(
		endpoints.FindFirstTenant,
		util.DecodeHTTPGenericRequest[tenantendpoint.FindFirstTenantRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	getAllTenantsHandler := httptransport.NewServer(
		endpoints.GetAllTenants,
		util.DecodeHTTPGenericRequest[tenantendpoint.GetAllTenantsRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	getByIdHandler := httptransport.NewServer(
		endpoints.GetTenantById,
		util.DecodeGenericHTTPIdentifierRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	findTenantsHandler := httptransport.NewServer(
		endpoints.FindTenants,
		util.DecodeHTTPGenericRequest[tenantendpoint.FindTenantsRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	updateTenantHandler := httptransport.NewServer(
		endpoints.UpdateTenant,
		util.DecodeHTTPGenericRequest[tenantendpoint.UpdateTenantRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	pagedTenantsHandler := httptransport.NewServer(
		endpoints.PagedTenants,
		util.DecodeHTTPPagedRequest[tenantendpoint.FindTenantsRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	deleteTenantHandler := httptransport.NewServer(
		endpoints.DeleteTenant,
		util.DecodeGenericHTTPIdentifierRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	return append(tenantServers, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: createHandler,
		Route: "/tenant",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: findFirstHandler,
		Route: "/tenant",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: getAllTenantsHandler,
		Route: "/tenant/all",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: findTenantsHandler,
		Route: "/tenant/find",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: getByIdHandler,
		Route: "/tenant/{identifier}",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: updateTenantHandler,
		Route: "/tenant",
		Methods: []string{"PUT"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: pagedTenantsHandler,
		Route: "/tenant/paged/",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: deleteTenantHandler,
		Route: "/tenant/{identifier}",
		Methods: []string{"DELETE"},
	})
}