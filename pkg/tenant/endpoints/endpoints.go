package tenantendpoint

import (
	"context"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/util"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	CreateTenant endpoint.Endpoint
}

func New(svc repository.BaseCRUD[tenantmodels.Tenant], logger log.Logger) Set {
	var createTenantEndpoint endpoint.Endpoint
	{
		createTenantEndpoint = CreateTenantEndpoint(svc)
		createTenantEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "createTenant"))(createTenantEndpoint)
	}

	return Set{
		CreateTenant: createTenantEndpoint,
	}
}

func CreateTenantEndpoint(s repository.BaseCRUD[tenantmodels.Tenant]) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (interface{}, error){
		req := request.(CreateTenantRequest)

		resp, err := s.Create(req.NewTenant)

		if err != nil {
			return CreateTenantResponse{resp, err}, err
		}

		return CreateTenantResponse{resp, nil}, nil

	}
}

type CreateTenantRequest struct {
	NewTenant tenantmodels.Tenant `json:"newTenant"`
}

type CreateTenantResponse struct {
	NewTenant tenantmodels.Tenant `json:"newTenant"`
	Err error `json:"err"`
}