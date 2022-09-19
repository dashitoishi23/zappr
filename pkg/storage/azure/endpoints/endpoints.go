package azurestorageendpoints

import (
	"context"

	azurestorageservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/storage/azure/service"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	CreateContainer endpoint.Endpoint
	DeleteContainer endpoint.Endpoint
}

func New(svc azurestorageservice.AzureStorageService, logger log.Logger) Set {
	var createContainerEndpoint endpoint.Endpoint 
	{
		createContainerEndpoint = CreateContainerEndpoint(svc)
		createContainerEndpoint = util.
		TransportLoggingMiddleware(log.With(logger, "method", "createContainer"))(createContainerEndpoint)
	}

	var deleteContainerEndpoint endpoint.Endpoint
	{
		deleteContainerEndpoint = DeleteContainerEndpoint(svc)
		deleteContainerEndpoint = util.
		TransportLoggingMiddleware(log.With(logger, "method", "deleteContainer"))(deleteContainerEndpoint)
	}

	return Set{
		CreateContainer: createContainerEndpoint,
		DeleteContainer: deleteContainerEndpoint,
	}
}

func CreateContainerEndpoint(a azurestorageservice.AzureStorageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateContainerRequest)

		isCreated, err := a.CreateAzureStorageContainer(ctx, req.ContainerName, req.IsPubliclyAccessible)

		return CreateContainerResponse{isCreated, err}, err
	}
}

func DeleteContainerEndpoint(a azurestorageservice.AzureStorageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteContainerRequest)

		isDeleted, err := a.DeleteAzureStorageContainer(ctx, req.ContainerName)

		return DeleteContainerResponse{isDeleted, err}, err
	}
}