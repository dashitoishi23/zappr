package initsetuptransports

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	initsetupendpoints "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/initsetup/endpoints"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(endpoints initsetupendpoints.Set) []commonmodels.HttpServerConfig {
	var configServers []commonmodels.HttpServerConfig

	serverOptions := []httptransport.ServerOption{
	httptransport.ServerErrorEncoder(util.ErrorEncoder),
	}

	addConfigHandler := httptransport.NewServer(
		endpoints.AddConfig,
		util.DecodeHTTPGenericRequest[initsetupendpoints.AddConfigRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	editConfigHandler := httptransport.NewServer(
		endpoints.EditConfig,
		util.DecodeHTTPPagedRequest[initsetupendpoints.EditConfigRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	configServers = append(configServers, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: addConfigHandler,
		Route: "/config/add",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: false,
		Server: editConfigHandler,
		Route: "/config/edit",
		Methods: []string{"PUT"},
	})

	return configServers

}