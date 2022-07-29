package masterroletransports

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	masterroleendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/role/endpoints"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func NewHandler(endpoints masterroleendpoint.Set, logger log.Logger) []commonmodels.HttpServerConfig {
	var roleServers []commonmodels.HttpServerConfig

	serverOptions := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(util.ErrorEncoder),
	}

	createHandler := httptransport.NewServer(
		endpoints.CreateRole,
		util.DecodeHTTPGenericRequest[masterroleendpoint.CreateRoleRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...,
	)

	findFirstHandler := httptransport.NewServer(
		endpoints.FindFirstRole,
		util.DecodeHTTPGenericRequest[masterroleendpoint.FindFirstRoleRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...,
	)

	getAllRolesHandler := httptransport.NewServer(
		endpoints.GetAllRoles,
		util.DecodeHTTPGenericRequest[masterroleendpoint.GetAllRolesRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...,
	)

	findRolesHandler := httptransport.NewServer(
		endpoints.FindRoles,
		util.DecodeHTTPGenericRequest[masterroleendpoint.FindRolesRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...,
	)

	updateRoleHandler := httptransport.NewServer(
		endpoints.UpdateRole,
		util.DecodeHTTPGenericRequest[masterroleendpoint.UpdateRoleRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...,
	)

	pagedRolesHandler := httptransport.NewServer(
		endpoints.PagedRolesEndpoint,
		util.DecodeHTTPPagedRequest[masterroleendpoint.FindRolesRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...,
	)

	return append(roleServers, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server:    createHandler,
		Route:     "/role",
		Methods:   []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server:    findFirstHandler,
	Route:     "/role",
		Methods:   []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server:    getAllRolesHandler,
		Route:     "/role/all",
		Methods:   []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server:    findRolesHandler,
		Route:     "/role/find",
		Methods:   []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server:    updateRoleHandler,
		Route:     "/role",
		Methods:   []string{"PUT"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server:    pagedRolesHandler,
		Route:     "/role/paged/",
		Methods:   []string{"GET"},
	})
}