package usermetadataendpoints

import (
	"context"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"

	usermetadatamodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/usermetadata/models"
)

type Set struct {
	AddUserMetadata endpoint.Endpoint
}

func New(svc repository.BaseCRUD[usermetadatamodels.UserMetadata], logger log.Logger) Set {
	var addUserMetadataEndpoint endpoint.Endpoint
	{
		addUserMetadataEndpoint = AddUserMetadataEndpoint(svc)
		addUserMetadataEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "findFirstTenant"))(addUserMetadataEndpoint)
	}

	return Set{
		AddUserMetadata: addUserMetadataEndpoint,
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