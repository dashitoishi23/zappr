package tenantendpoint

import (
	"context"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	tenantservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/service"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	CreateTenant endpoint.Endpoint
	FindFirstTenant endpoint.Endpoint
	GetAllTenants endpoint.Endpoint
	GetTenantById endpoint.Endpoint
	FindTenants endpoint.Endpoint
	UpdateTenant endpoint.Endpoint
	PagedTenants endpoint.Endpoint
	DeleteTenant endpoint.Endpoint
}

func New(svc repository.BaseCRUD[tenantmodels.Tenant], tenantSvc tenantservice.TenantService, logger log.Logger) Set {
	var createTenantEndpoint endpoint.Endpoint
	{
		createTenantEndpoint = CreateTenantEndpoint(tenantSvc)
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

	var getTenantByIdEndpoint endpoint.Endpoint
	{
		getTenantByIdEndpoint = GetTenantByIdEndpoint(svc)
		getTenantByIdEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "getTenantById"))(getTenantByIdEndpoint)
	}

	var updatedTenantEndpoint endpoint.Endpoint
	{
		updatedTenantEndpoint = UpdateTenantEndpoint(svc)
		updatedTenantEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "updateTenant"))(updatedTenantEndpoint)
	}

	var getPagedTenantsEndpoint endpoint.Endpoint
	{
		getPagedTenantsEndpoint = GetPagedTenantsEndpoint(svc)
		getPagedTenantsEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "pagedTenantsEndpoint"))(getPagedTenantsEndpoint)
	}

	var deleteTenantEndpoint endpoint.Endpoint
	{
		deleteTenantEndpoint = DeleteTenantEndpoint(svc)
		deleteTenantEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "deleteTenant"))(deleteTenantEndpoint)
	}

	return Set{
		CreateTenant: createTenantEndpoint,
		FindFirstTenant: findFirstTenantEndpoint,
		GetAllTenants: getAllTenantsEndpoint,
		GetTenantById: getTenantByIdEndpoint,
		FindTenants: findTenantsEndpoint,
		UpdateTenant: updatedTenantEndpoint,
		PagedTenants: getPagedTenantsEndpoint,
		DeleteTenant: deleteTenantEndpoint,
	}
}

func CreateTenantEndpoint(s tenantservice.TenantService) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (interface{}, error){
		req := request.(CreateTenantRequest)

		req.NewTenant.InitFields()

		resp, err := s.AddTenant(req.NewTenant)

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

func GetTenantByIdEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp tenantmodels.Tenant

		req := request.(string)

		resp, err = s.GetFirst(&tenantmodels.Tenant{
			Identifier: req,
		})

		return resp, err
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

func UpdateTenantEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
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
				close(currentEntity)
				close(txnError)
				return UpdateTenantResponse{resp, err}, err
			}	
		}
	}
}

func GetPagedTenantsEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(commonmodels.PagedRequest[FindTenantsRequest])

		res, err := s.GetPaged(req.Entity.CurrentTenant, req.Page, req.Size)

		return GetPagedTenantResponse{res, err}, err
	}
}

func DeleteTenantEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		resp, err := s.Delete(&tenantmodels.Tenant{
			Identifier: req,
		})

		return resp, err
	}
}

