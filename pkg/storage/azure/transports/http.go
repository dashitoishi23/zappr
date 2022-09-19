package azurestoragetransports

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	azurestorageendpoints "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/storage/azure/endpoints"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(endpoints azurestorageendpoints.Set) []commonmodels.HttpServerConfig {
	var azureEndpointServers []commonmodels.HttpServerConfig

	serverOptions := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(util.ErrorEncoder),
	}

	createContainerHandler := httptransport.NewServer(
		endpoints.CreateContainer,
		util.DecodeHTTPGenericRequest[azurestorageendpoints.CreateContainerRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	deleteContainerHandler := httptransport.NewServer(
		endpoints.DeleteContainer,
		util.DecodeHTTPGenericRequest[azurestorageendpoints.DeleteContainerRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	return append(azureEndpointServers, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: createContainerHandler,
		Route: "/storage/azure/container",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: deleteContainerHandler,
		Route: "/storage/azure/container",
		Methods: []string{"DELETE"},
	})
}