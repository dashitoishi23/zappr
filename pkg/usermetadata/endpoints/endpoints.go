package usermetadataendpoints

import (
	"context"
	"encoding/json"
	"errors"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/state"
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

	return Set{
		AddUserMetadata: addUserMetadataEndpoint,
		GetUserMetadata: getUserMetadataEndpoint,
		GetMetadataByEntity: getMetadataByEndpoint,
		GetMetadataByEntityPaged: getMetadataByEntityPagedEndpoint,
	}
}

func AddUserMetadataEndpoint(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddUserMetadataRequest)

		req.NewUserMetadata.InitFields()

		res, err := s.Create(req.NewUserMetadata)

		return AddUserMetadataResponse{res, err}, err
	}
}

func GetUserMetadataEndpoint(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserMetadataRequest)

		query, queryErr := constructJSONBQuery(req.Query, req.EntityName)

		if queryErr != nil {
			return nil, queryErr
		}

		res, err := s.QueryRawSql(query)

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
			TenantIdentifier: state.GetState().UserContext.UserTenant,
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

		query, queryErr := constructJSONBQuery(req.Query, req.EntityName)

		if queryErr != nil {
			return nil, queryErr
		}

		res, err := s.QueryRawSqlPaged(query, req.Page, req.Size)

		return res, err
		
	}
}

func constructJSONBQuery(query map[string]interface{}, entityName string) (string, error) {
	jsonQuery, jsonErr := json.Marshal(query)

		if jsonErr != nil {
			return "", jsonErr
		}

		return "select * from \"UserMetadata\" where \"Metadata\" @> '" + string(jsonQuery) + "' and \"TenantIdentifier\" = '" + state.GetState().UserContext.UserTenant + "' and \"EntityName\" = '"  + entityName + "' ", nil
}