package tenantendpoint

import (
	"context"
	"fmt"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	tenantservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/service"
	redisutil "dev.azure.com/technovert-vso/Zappr/_git/Zappr/redis"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/gomodule/redigo/redis"
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

func New(svc repository.BaseCRUD[tenantmodels.Tenant], tenantSvc tenantservice.TenantService, logger log.Logger, 
	client redis.Conn) Set {
	var createTenantEndpoint endpoint.Endpoint
	{
		createTenantEndpoint = CreateTenantEndpoint(tenantSvc)
		createTenantEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "createTenant"))(createTenantEndpoint)
	}

	var findFirstTenantEndpoint endpoint.Endpoint
	{
		findFirstTenantEndpoint = FindFirstTenantEndpoint(svc, client)
		findFirstTenantEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "findFirstTenant"))(findFirstTenantEndpoint)
	}

	var getAllTenantsEndpoint endpoint.Endpoint
	{
		getAllTenantsEndpoint = GetAllTenantsEndpoint(svc, client)
		getAllTenantsEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "getAllTenants"))(getAllTenantsEndpoint)
	}

	var findTenantsEndpoint endpoint.Endpoint
	{
		findTenantsEndpoint = FindTenantsEndpoint(svc)
		findTenantsEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "findTenants"))(findTenantsEndpoint)
	}

	var getTenantByIdEndpoint endpoint.Endpoint
	{
		getTenantByIdEndpoint = GetTenantByIdEndpoint(svc, client)
		getTenantByIdEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "getTenantById"))(getTenantByIdEndpoint)
	}

	var updatedTenantEndpoint endpoint.Endpoint
	{
		updatedTenantEndpoint = UpdateTenantEndpoint(svc, client)
		updatedTenantEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "updateTenant"))(updatedTenantEndpoint)
	}

	var getPagedTenantsEndpoint endpoint.Endpoint
	{
		getPagedTenantsEndpoint = GetPagedTenantsEndpoint(svc, client)
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

func FindFirstTenantEndpoint(s repository.BaseCRUD[tenantmodels.Tenant], client redis.Conn) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp tenantmodels.Tenant

		req := request.(FindFirstTenantRequest)

		stringifiedQuery, err := util.StringifyJson(req.CurrentTenant)

		if err != nil {
			return nil, err
		}

		cacheKey := "tenant:findFirstTenant:" + stringifiedQuery

		cachedResponse, err := redisutil.DecodeCacheResponse[tenantmodels.Tenant](client, cacheKey)

		if err != nil {
			return nil, err
		}

		if len(cachedResponse.Identifier) != 0 {
			return FindFirstTenantResponse{cachedResponse, nil}, nil
		}

		resp, err = s.GetFirst(req.CurrentTenant)

		if err != nil {
			return nil, err
		}

		err = redisutil.SetCache(client, cacheKey, resp)

		return FindFirstTenantResponse{resp, err}, err
	}
}

func GetAllTenantsEndpoint(s repository.BaseCRUD[tenantmodels.Tenant], client redis.Conn) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []tenantmodels.Tenant

		cacheKey := "tenant:getAll" 

		decodedResponse, err := redisutil.DecodeCacheResponse[[]tenantmodels.Tenant](client, cacheKey)

		if err != nil {
			return nil, err
		}

		if decodedResponse != nil {
			return GetAllTenantsResponse{decodedResponse, nil}, nil
		}

		resp, err = s.GetAll()

		if err != nil {
			return nil, err
		}

		err = redisutil.SetCache(client, cacheKey, resp)

		return GetAllTenantsResponse{resp, err}, err

	}
}

func GetTenantByIdEndpoint(s repository.BaseCRUD[tenantmodels.Tenant], client redis.Conn) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp tenantmodels.Tenant

		req := request.(string)

		cacheKey := "tenant:getTenantById:" + req

		cachedResponse, err := redisutil.DecodeCacheResponse[tenantmodels.Tenant](client, cacheKey)

		if err != nil {
			return nil, err
		}

		fmt.Print(cachedResponse)

		if len(cachedResponse.Identifier) != 0 {
			return cachedResponse, nil
		}

		resp, err = s.GetFirst(&tenantmodels.Tenant{
			Identifier: req,
		})

		if err != nil {
			return nil, err
		}

		err = redisutil.SetCache(client, cacheKey, resp)

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

func UpdateTenantEndpoint(s repository.BaseCRUD[tenantmodels.Tenant], client redis.Conn) endpoint.Endpoint {
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

				if err != nil {
					return nil, err
				}

				tenantKeys, err := client.Do("KEYS", "tenant:*")

				if err != nil {
					return nil, err
				}

				keys := util.StringifyTo2dArray(tenantKeys.([]interface{}))

				delErr := redisutil.DeleteMultipleKeys(client, keys)

				if delErr != nil {
					return nil, delErr
				}

				return UpdateTenantResponse{resp, err}, err
			}	
		}
	}
}

func GetPagedTenantsEndpoint(s repository.BaseCRUD[tenantmodels.Tenant], client redis.Conn) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(commonmodels.PagedRequest[FindTenantsRequest])

		stringifiedQuery, err := util.StringifyJson(req)

		if err != nil {
			return nil, err
		}

		cacheKey := "tenant:getPagedTenants" + stringifiedQuery

		cachedResponse, err := redisutil.DecodeCacheResponse[commonmodels.PagedResponse[tenantmodels.Tenant]](client, cacheKey)

		if err != nil {
			return nil, err
		}

		if cachedResponse.Pages != 0 {
			return GetPagedTenantResponse{cachedResponse, nil}, nil
		}

		res, err := s.GetPaged(req.Entity.CurrentTenant, req.Page, req.Size)

		if err != nil {
			return nil, err
		}

		err = redisutil.SetCache(client, cacheKey, res)

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

