package tenanttransports

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	tenantendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/endpoints"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHandler(endpoints tenantendpoint.Set) []commonmodels.HttpServerConfig {
	var tenantServers []commonmodels.HttpServerConfig

	createHandler := httptransport.NewServer(
		endpoints.CreateTenant,
		decodeCreateTenantRequest,
		util.EncodeHTTPGenericResponse,
	)

	findFirstHandler := httptransport.NewServer(
		endpoints.FindFirstTenant,
		decodeFindFirstTenantRequest,
		util.EncodeHTTPGenericResponse,
	)

	getAllTenantsHandler := httptransport.NewServer(
		endpoints.GetAllTenants,
		decodeGetAllTenantsRequest,
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
	})
}

func decodeCreateTenantRequest(_ context.Context, r *http.Request) (interface{}, error){
	var req tenantendpoint.CreateTenantRequest

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	if err != nil {
		return req, err
	}

	return req, nil

}

func decodeFindFirstTenantRequest(_ context.Context, r *http.Request) (interface{}, error){
	var req tenantendpoint.FindFirstTenantRequest

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	if err != nil {
		return req, err
	}

	return req, nil

}

func decodeGetAllTenantsRequest(_ context.Context, r *http.Request) (interface{}, error){
	var req tenantendpoint.GetAllTenantsRequest

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	if err == io.EOF {
		return req, nil
	}

	if err != nil {
		return req, err
	}

	return req, nil

}