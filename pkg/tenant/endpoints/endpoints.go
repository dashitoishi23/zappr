package tenantendpoint

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	CreateTenant endpoint.Endpoint
	FindFirstTenant endpoint.Endpoint
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

	return Set{
		CreateTenant: createTenantEndpoint,
		FindFirstTenant: findFirstTenantEndpoint,
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
		var reqBody tenantmodels.SearchableTenant

		req := request.(FindFirstTenantRequest)

		reqData := json.NewDecoder(strings.NewReader(request.(string)))

		reqData.DisallowUnknownFields()

		if errs := reqData.Decode(&reqBody); errs != nil {
			return resp, errors.New(constants.INVALID_MODEL)
		}

		resp, err = s.GetFirst(req.CurrentTenant)

		if err != nil {
			return FindFirstTenantResponse{resp}, err
		}

		return FindFirstTenantResponse{resp}, nil
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