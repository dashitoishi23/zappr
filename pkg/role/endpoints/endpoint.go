package masterroleendpoint

import (
	"context"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	masterrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/role/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	CreateRole         endpoint.Endpoint
	FindFirstRole      endpoint.Endpoint
	GetAllRoles        endpoint.Endpoint
	FindRoles          endpoint.Endpoint
	UpdateRole         endpoint.Endpoint
	PagedRolesEndpoint endpoint.Endpoint
}

func New(svc repository.BaseCRUD[masterrolemodels.Role], logger log.Logger) Set {
	var createRoleEndpoint endpoint.Endpoint
	{
		createRoleEndpoint = CreateRoleEndpoint(svc)
		createRoleEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "createRole"))(createRoleEndpoint)
	}

	var findFirstRoleEndpoint endpoint.Endpoint
	{
		findFirstRoleEndpoint = FindFirstRoleEndpoint(svc)
		findFirstRoleEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "findFirstRole"))(findFirstRoleEndpoint)
	}

	var getAllRolesEndpoint endpoint.Endpoint
	{
		getAllRolesEndpoint = GetAllRolesEndpoint(svc)
		getAllRolesEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "getAllRoles"))(getAllRolesEndpoint)
	}

	var findRolesEndpoint endpoint.Endpoint
	{
		findRolesEndpoint = FindRolesEndpoint(svc)
		findRolesEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "findRoles"))(findRolesEndpoint)
	}

	var updatedRoleEndpoint endpoint.Endpoint
	{
		updatedRoleEndpoint = UpdateRoleEndpoint(svc)
		updatedRoleEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "updateRole"))(updatedRoleEndpoint)
	}

	var getPagedRolesEndpoint endpoint.Endpoint
	{
		getPagedRolesEndpoint = GetPagedRolesEndpoint(svc)
		getPagedRolesEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "pagedRolesEndpoint"))(getPagedRolesEndpoint)
	}

	return Set{
		CreateRole: createRoleEndpoint,
		FindFirstRole: findFirstRoleEndpoint,
		GetAllRoles: getAllRolesEndpoint,
		FindRoles: findRolesEndpoint,
		UpdateRole: updatedRoleEndpoint,
		PagedRolesEndpoint: getPagedRolesEndpoint,
	}
}

func CreateRoleEndpoint(s repository.BaseCRUD[masterrolemodels.Role]) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (interface{}, error){
		req := request.(CreateRoleRequest)

		req.NewRole.InitFields(ctx.Value("requestScope").(commonmodels.RequestScope))

		resp, err := s.Create(req.NewRole)

		if err != nil {
			return CreateRoleResponse{resp, err}, err
		}

		return CreateRoleResponse{resp, err}, nil

	}
}

func FindFirstRoleEndpoint(s repository.BaseCRUD[masterrolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp masterrolemodels.Role

		req := request.(FindFirstRoleRequest)

		req.CurrentRole.AddTenant(ctx.Value("requestScope").(commonmodels.RequestScope))

		resp, err = s.GetFirst(req.CurrentRole)

		if err != nil {
			return FindFirstRoleResponse{resp, err}, err
		}

		return FindFirstRoleResponse{resp, err}, nil
	}
}

func GetAllRolesEndpoint(s repository.BaseCRUD[masterrolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []masterrolemodels.Role

		resp, err = s.GetAllByTenant(ctx)

		if err != nil {
			return GetAllRolesResponse{resp, err}, err
		}

		return GetAllRolesResponse{resp, err}, nil

	}
}

func FindRolesEndpoint(s repository.BaseCRUD[masterrolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []masterrolemodels.Role

		req := request.(FindRolesRequest)

		req.CurrentRole.AddTenant(ctx.Value("requestScope").(commonmodels.RequestScope))

		resp, err = s.Find(req.CurrentRole)

		if err != nil {
			return FindRolesResponse{resp, err}, err
		}

		return FindRolesResponse{resp, err}, nil

	}
}

func UpdateRoleEndpoint(s repository.BaseCRUD[masterrolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateRoleRequest)

		currentEntity := make(chan masterrolemodels.Role)
		txnError := make(chan error)

		go s.GetFirstAsync(&masterrolemodels.Role{
			Identifier: req.NewRole.Identifier,
		}, currentEntity, txnError)

		for{
			select {			
			case entity := <-currentEntity:
				err := <-txnError
				req.NewRole.UpdateFields(entity.CreatedOn, ctx.Value("requestScope").(commonmodels.RequestScope))
				resp, _ := s.Update(req.NewRole)				
				close(currentEntity)
				close(txnError)
				return UpdateRoleResponse{resp, err}, err
			}	
		}
	}
}

func GetPagedRolesEndpoint(s repository.BaseCRUD[masterrolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(commonmodels.PagedRequest[FindRolesRequest])

		respChan := make(chan commonmodels.PagedResponse[masterrolemodels.Role])
		txnError := make(chan error)

		req.Entity.CurrentRole.AddTenant(ctx.Value("requestScope").(commonmodels.RequestScope))
		go s.GetPagedAsync(req.Entity.CurrentRole, req.Page, req.Size, respChan, txnError)

		for{
			select {			
			case entity := <-respChan:
				err := <-txnError			
				close(respChan)
				close(txnError)
				return commonmodels.PagedResponse[masterrolemodels.Role]{
					Items: entity.Items,
					Page: entity.Page,
					Size: entity.Size,
					Pages: entity.Pages,
					Err: err,
				}, err
			}	
		}
	}
}
