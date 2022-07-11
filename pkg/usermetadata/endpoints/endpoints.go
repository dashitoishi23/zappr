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

	return Set{
		AddUserMetadata: addUserMetadataEndpoint,
		GetUserMetadata: getUserMetadataEndpoint,
	}
}

func AddUserMetadataEndpoint(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddUserMetadataRequest)

		req.NewUserMetadata.InitFields()

		res, err := s.Create(req.NewUserMetadata)

		return res, err
	}
}

func GetUserMetadataEndpoint(s repository.BaseCRUD[usermetadatamodels.UserMetadata]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserMetadataRequest)

		jsonQuery, jsonErr := json.Marshal(req.Query)

		if jsonErr != nil {
			return nil, jsonErr
		}

		query := "select * from \"UserMetadata\" where \"Metadata\" @> '" + string(jsonQuery) + "' and \"TenantIdentifier\" = '" + state.GetState().UserContext.UserTenant + "' and \"EntityName\" = '"  + req.EntityName + "' "

		res, err := s.QueryRawSql(query)

		if len(res) == 0 {
			return res, errors.New(constants.RECORD_NOT_FOUND)
		}

		return res, err
	}
}