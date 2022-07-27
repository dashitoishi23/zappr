package usermetadataendpoints

import (
	"context"
	"encoding/json"
	"errors"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"

	usermetadatamodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/usermetadata/models"
)

type Set struct {
	AddUserMetadata endpoint.Endpoint
	GetUserMetadata endpoint.Endpoint
	GetMetadataByEntity endpoint.Endpoint
	GetMetadataByEntityPaged endpoint.Endpoint
	UpdateMetadata endpoint.Endpoint
	DeleteMetadata endpoint.Endpoint
}

func New(svc repository.BaseCRUD[usermetadatamodels.UserMetadata], logger log.Logger) Set {
	var addUserMetadataEndpoint endpoint.Endpoint
	{
		addUserMetadataEndpoint = AddUserMetadataEndpoint(svc)
		addUserMetadataEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "findFirstTenant"))(addUserMetadataEndpoint)
	}

	var getUserMetadataEndpoint endpoint.Endpoint
	{
		getUserMetadataEndpoint = GetUserMetadataEndpoint(svc)
		getUserMetadataEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "getUserMetadata"))(getUserMetadataEndpoint)
	}

	var getMetadataByEndpoint endpoint.Endpoint
	{
		getMetadataByEndpoint = GetMetadataByEntity(svc)
		getMetadataByEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "getMetadataByEndpoint"))(getMetadataByEndpoint)
	}

	var getMetadataByEntityPagedEndpoint endpoint.Endpoint
	{
		getMetadataByEntityPagedEndpoint = GetMetadataByEntityPagedEndpoint(svc)
		getMetadataByEntityPagedEndpoint = util.TransportLoggingMiddleware(log.
			With(logger, "method", "getMetadataByEntityPagedEndpoint"))(getMetadataByEntityPagedEndpoint)
	}

	var updateMetdataEndpoint endpoint.Endpoint
	{
		updateMetdataEndpoint = UpdateMetadataEndpoint(svc)
		updateMetdataEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "updateMetadataEndpoint"))(updateMetdataEndpoint)
	}

	var deleteMetadataEndpoint endpoint.Endpoint
	{
		deleteMetadataEndpoint = DeleteMetadataEndpoint(svc)
		deleteMetadataEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "deleteMetadataEndpoint"))(deleteMetadataEndpoint)
	}

	return Set{
		AddUserMetadata: addUserMetadataEndpoint,
		GetUserMetadata: getUserMetadataEndpoint,
		GetMetadataByEntity: getMetadataByEndpoint,
		GetMetadataByEntityPaged: getMetadataByEntityPagedEndpoint,
		UpdateMetadata: updateMetdataEndpoint,
		DeleteMetadata: deleteMetadataEndpoint,
	}
}

func AddUserMetadataEndpoint(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddUserMetadataRequest)

		requestScope := ctx.Value("requestScope").(commonmodels.RequestScope)

		req.NewUserMetadata.InitFields(requestScope)

		res, err := s.Create(req.NewUserMetadata)

		return AddUserMetadataResponse{res, err}, err
	}
}

func GetUserMetadataEndpoint(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserMetadataRequest)

		tenantId := ctx.Value("requestScope").(commonmodels.RequestScope).UserTenant

		query, queryErr := constructJSONBQuery(req.Query, req.EntityName, tenantId)

		if queryErr != nil {
			return nil, queryErr
		}

		res, err := s.QueryRawSql("select * from \"UserMetadata\" where \"Metadata\" @> '" + query)

		if len(res) == 0 {
			return res, errors.New(constants.RECORD_NOT_FOUND)
		}

		return GetUserMetadataResponse{res, err}, err
	}
}

func GetMetadataByEntity(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		entityName := request.(string)

		res, err := s.Find(usermetadatamodels.SearchableUserMetadata{
			TenantIdentifier: ctx.Value("requestScope").(commonmodels.RequestScope).UserTenant,
			EntityName: entityName,
		})

		var metadata []json.RawMessage

		if len(res) == 0 {
			return GetMetadataByEntityResponse{metadata, errors.New(constants.RECORD_NOT_FOUND)}, 
			errors.New(constants.RECORD_NOT_FOUND)
		}

		
		if err != nil {
			return GetMetadataByEntityResponse{metadata, err}, err
		}


		for _, data := range res {
			metadata = append(metadata, data.Metadata)
		}

		return GetMetadataByEntityResponse{metadata, nil}, nil

	}
}

func GetMetadataByEntityPagedEndpoint(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetMetadataByEntityPagedRequest)

		tenantId := ctx.Value("requestScope").(commonmodels.RequestScope).UserTenant

		query, queryErr := constructJSONBQuery(req.Query, req.EntityName, tenantId)

		if queryErr != nil {
			return nil, queryErr
		}

		res, err := s.QueryRawSqlPaged("select * from \"UserMetadata\" where \"Metadata\" @> '" + query, req.Page, req.Size)

		var pagedResponse commonmodels.PagedResponse[json.RawMessage]

		for _, item := range res.Items {
			pagedResponse.Items = append(pagedResponse.Items, item.Metadata)
		}

		pagedResponse.Page = res.Page
		pagedResponse.Size = res.Size
		pagedResponse.Pages = res.Pages

		return GetMetadataByEntityPagedResponse{pagedResponse, err}, err
		
	}
}

func UpdateMetadataEndpoint(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateMetadataRequest)

		requestScope :=  ctx.Value("requestScope").(commonmodels.RequestScope)

		tenantId := requestScope.UserTenant

		updatedQuery, queryErr := json.Marshal(req.UpdatedQuery)

		if queryErr != nil {
			return nil, queryErr
		}

		sql, queryErr := constructJSONBQuery(req.CurrentQuery, req.EntityName, tenantId)

		if queryErr != nil {
			return nil, queryErr
		}

		query := "update \"UserMetadata\" set \"Metadata\" = '" + string(updatedQuery) + "',\"ModifiedOn\" = NOW() where \"Metadata\" @> '" +  sql
		
		resp, err := s.ExecuteRawQuery(query)

		if !resp {
			return resp, err
		}

		return UpdateMetadataResponse{req.UpdatedQuery, nil}, nil
	}
}

func DeleteMetadataEndpoint(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserMetadataRequest)

		requestScope := ctx.Value("requestScope").(commonmodels.RequestScope)

		tenantId := requestScope.UserTenant

		query, queryErr := constructJSONBQuery(req.Query, req.EntityName, tenantId)

		if queryErr != nil {
			return nil, queryErr
		}

		sqlString := "delete from \"UserMetadata\" where \"Metadata\" @> '" + query

		resp, err := s.ExecuteRawQuery(sqlString)

		if !resp {
			return resp, err
		}

		return DeleteMetadataResponse{resp, nil}, nil

	}
}

func constructJSONBQuery(query map[string]interface{}, entityName string, tenantIdentifier string) (string, error) {
	jsonQuery, jsonErr := json.Marshal(query)

		if jsonErr != nil {
			return "", jsonErr
		}

		return string(jsonQuery) + "' and \"TenantIdentifier\" = '" + tenantIdentifier + "' and \"EntityName\" = '"  + entityName + "' ", nil
}