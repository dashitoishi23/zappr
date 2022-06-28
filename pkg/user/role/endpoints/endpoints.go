package rolesendpoint

import (
	"context"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	rolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/role/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	CreateRole endpoint.Endpoint
	FindFirstRole endpoint.Endpoint
	GetAllRoles endpoint.Endpoint
	FindRoles endpoint.Endpoint
	UpdateRole endpoint.Endpoint
	PagedRolesEndpoint endpoint.Endpoint
}

func New(svc repository.BaseCRUD[rolemodels.Role], logger log.Logger) Set {
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

func CreateRoleEndpoint(s repository.BaseCRUD[rolemodels.Role]) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (interface{}, error){
		req := request.(CreateRoleRequest)

		req.NewRole.InitFields()

		resp, err := s.Create(req.NewRole)

		if err != nil {
			return CreateRoleResponse{resp, err}, err
		}

		return CreateRoleResponse{resp, err}, nil

	}
}

func FindFirstRoleEndpoint(s repository.BaseCRUD[rolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp rolemodels.Role

		req := request.(FindFirstRoleRequest)

		resp, err = s.GetFirst(req.CurrentRole)

		if err != nil {
			return FindFirstRoleResponse{resp, err}, err
		}

		return FindFirstRoleResponse{resp, err}, nil
	}
}

func GetAllRolesEndpoint(s repository.BaseCRUD[rolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []rolemodels.Role

		resp, err = s.GetAll()

		if err != nil {
			return GetAllRolesResponse{resp, err}, err
		}

		return GetAllRolesResponse{resp, err}, nil

	}
}

func FindRolesEndpoint(s repository.BaseCRUD[rolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []rolemodels.Role

		req := request.(FindRolesRequest)

		resp, err = s.Find(req.CurrentRole)

		if err != nil {
			return FindRolesResponse{resp, err}, err
		}

		return FindRolesResponse{resp, err}, nil

	}
}

func UpdateRoleEndpoint(s repository.BaseCRUD[rolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateRoleRequest)

		currentEntity := make(chan rolemodels.Role)
		txnError := make(chan error)

		go s.GetFirstAsync(&rolemodels.Role{
			Identifier: req.NewRole.Identifier,
		}, currentEntity, txnError)

		for{
			select {			
			case entity := <-currentEntity:
				err := <-txnError
				req.NewRole.UpdateFields(entity.CreatedOn)
				resp, _ := s.Update(req.NewRole)				
				close(currentEntity)
				close(txnError)
				return UpdateRoleResponse{resp, err}, err
			}	
		}
	}
}

func GetPagedRolesEndpoint(s repository.BaseCRUD[rolemodels.Role]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(commonmodels.PagedRequest[FindRolesRequest])

		respChan := make(chan commonmodels.PagedResponse[rolemodels.Role])
		txnError := make(chan error)
		go s.GetPagedAsync(req.Entity.CurrentRole, req.Page, req.Size, respChan, txnError)

		for{
			select {			
			case entity := <-respChan:
				err := <-txnError			
				close(respChan)
				close(txnError)
				return commonmodels.PagedResponse[rolemodels.Role]{
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

