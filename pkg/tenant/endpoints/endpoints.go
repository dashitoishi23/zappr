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
	UpdateTenant endpoint.Endpoint
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

	var updatedTenantEndpoint endpoint.Endpoint
	{
		updatedTenantEndpoint = UpdateTenantEndpint(svc)
		updatedTenantEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "updateTenant"))(updatedTenantEndpoint)
	}

	return Set{
		CreateTenant: createTenantEndpoint,
		FindFirstTenant: findFirstTenantEndpoint,
		GetAllTenants: getAllTenantsEndpoint,
		FindTenants: findTenantsEndpoint,
		UpdateTenant: updatedTenantEndpoint,
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

func UpdateTenantEndpint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateTenantRequest)

		currentEntity := make(chan tenantmodels.Tenant)
		txnError := make(chan error)

		go s.GetFirstAsync(&tenantmodels.Tenant{
			Identifier: req.NewTenant.Identifier,
		}, currentEntity, txnError)

		for{
			select {			
			case entity := <-currentEntity:
				err := <-txnError
				req.NewTenant.UpdateFields(entity.CreatedOn)
				resp, _ := s.Update(req.NewTenant)
				return UpdateTenantResponse{resp, err}, err
			}
		}
	}
}

