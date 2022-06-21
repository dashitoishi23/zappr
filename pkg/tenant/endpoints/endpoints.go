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
			return CreateTenantResponse{resp, err}, err
		}

		return CreateTenantResponse{resp, err}, nil

	}
}

func FindFirstTenantEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp tenantmodels.Tenant

		req := request.(FindFirstTenantRequest)

		resp, err = s.GetFirst(req.CurrentTenant)

		if err != nil {
			return FindFirstTenantResponse{resp, err}, err
		}

		return FindFirstTenantResponse{resp, err}, nil
	}
}

func GetAllTenantsEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []tenantmodels.Tenant

		resp, err = s.GetAll()

		if err != nil {
			return GetAllTenantsResponse{resp, err}, err
		}

		return GetAllTenantsResponse{resp, err}, nil

	}
}

func FindTenantsEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []tenantmodels.Tenant

		req := request.(FindTenantsRequest)

		resp, err = s.Find(req.CurrentTenant)

		if err != nil {
			return FindTenantsResponse{resp, err}, err
		}

		return FindTenantsResponse{resp, err}, nil

	}
}

type CreateTenantRequest struct {
	NewTenant tenantmodels.Tenant `json:"newTenant"`
}

type CreateTenantResponse struct {
	NewTenant tenantmodels.Tenant `json:"newTenant"`
	Err error `json:"-"`
}

func (c CreateTenantResponse) Failed() error { return c.Err }

type FindFirstTenantRequest struct {
	CurrentTenant tenantmodels.SearchableTenant `json:"currentTenant"`
}

type FindFirstTenantResponse struct {
	CurrentTenant tenantmodels.Tenant `json:"currentTenant"`
	Err error `json:"-"`
}

func (f FindFirstTenantResponse) Failed() error { return f.Err }

type GetAllTenantsRequest struct {

}

type GetAllTenantsResponse struct {
	Tenants []tenantmodels.Tenant `json:"tenants"`
	Err error `json:"-"`
}

func (g GetAllTenantsResponse) Failed() error { return g.Err }

type FindTenantsRequest struct {
	CurrentTenant tenantmodels.SearchableTenant `json:"currentTenant"`
}

type FindTenantsResponse struct {
	CurrentTenant []tenantmodels.Tenant `json:"currentTenant"`
	Err error `json:"-"`
}

func (f FindTenantsResponse) Failed() error { return f.Err }