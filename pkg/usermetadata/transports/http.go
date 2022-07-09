package usermetadatatransports

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	usermetadataendpoints "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/usermetadata/endpoints"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func NewHandler(endpoints usermetadataendpoints.Set, logger log.Logger) []commonmodels.HttpServerConfig {
	var usermetadataServers []commonmodels.HttpServerConfig

	serverOptions := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(util.ErrorEncoder),
	}

	addHandler := httptransport.NewServer(
		endpoints.AddUserMetadata,
		util.DecodeHTTPGenericRequest[usermetadataendpoints.AddUserMetadataRequest],
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	return append(usermetadataServers, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: addHandler,
		Route:"/usermetadata",
		Methods: []string{"POST"},
	})
}