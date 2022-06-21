package tenantendpoint

import (
	"context"

	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	CreateTenant endpoint.Endpoint
	FindFirstTenant endpoint.Endpoint
	GetAllTenants endpoint.Endpoint
	FindTenants endpoint.Endpoint
}

func New(svc repository.BaseCRUD[tenantmodels.Tenant], logger log.Logger) Set {
	var createTenantEndpoint endpoint.Endpoint
	{
		createTenantEndpoint = CreateTenantEndpoint(svc)
		createTenantEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "createTenant"))(createTenantEndpoint)
	}

	var findFirstTenantEndpoint endpoint.Endpoint
	{
		findFirstTenantEndpoint = FindFirstTenantEndpoint(svc)
		findFirstTenantEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "findFirstTenant"))(findFirstTenantEndpoint)
	}

	var getAllTenantsEndpoint endpoint.Endpoint
	{
		getAllTenantsEndpoint = GetAllTenantsEndpoint(svc)
		getAllTenantsEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "getAllTenants"))(getAllTenantsEndpoint)
	}

	var findTenantsEndpoint endpoint.Endpoint
	{
		findTenantsEndpoint = FindTenantsEndpoint(svc)
		findTenantsEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "findTenants"))(findTenantsEndpoint)
	}

	return Set{
		CreateTenant: createTenantEndpoint,
		FindFirstTenant: findFirstTenantEndpoint,
		GetAllTenants: getAllTenantsEndpoint,
		FindTenants: findTenantsEndpoint,
	}
}

func CreateTenantEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (interface{}, error){
		req := request.(CreateTenantRequest)

		req.NewTenant.InitFields()

		resp, err := s.Create(req.NewTenant)

		if err != nil {
			return CreateTenantResponse{resp}, err
		}

		return CreateTenantResponse{resp}, nil

	}
}

func FindFirstTenantEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp tenantmodels.Tenant

		req := request.(FindFirstTenantRequest)

		resp, err = s.GetFirst(req.CurrentTenant)

		if err != nil {
			return FindFirstTenantResponse{resp}, err
		}

		return FindFirstTenantResponse{resp}, nil
	}
}

func GetAllTenantsEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []tenantmodels.Tenant

		resp, err = s.GetAll()

		if err != nil {
			return GetAllTenantsResponse{resp}, err
		}

		return GetAllTenantsResponse{resp}, nil

	}
}

func FindTenantsEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []tenantmodels.Tenant

		req := request.(FindTenantsRequest)

		resp, err = s.Find(req.CurrentTenant)

		if err != nil {
			return FindTenantsResponse{resp}, err
		}

		return FindTenantsResponse{resp}, nil

	}
}

type CreateTenantRequest struct {
	NewTenant tenantmodels.Tenant `json:"newTenant"`
}

type CreateTenantResponse struct {
	NewTenant tenantmodels.Tenant `json:"newTenant"`
}

type FindFirstTenantRequest struct {
	CurrentTenant tenantmodels.SearchableTenant `json:"currentTenant"`
}

type FindFirstTenantResponse struct {
	CurrentTenant tenantmodels.Tenant `json:"currentTenant"`
}

type GetAllTenantsRequest struct {

}

type GetAllTenantsResponse struct {
	Tenants []tenantmodels.Tenant `json:"tenants"`
}

type FindTenantsRequest struct {
	CurrentTenant tenantmodels.SearchableTenant `json:"currentTenant"`
}

type FindTenantsResponse struct {
	CurrentTenant []tenantmodels.Tenant `json:"currentTenant"`
}