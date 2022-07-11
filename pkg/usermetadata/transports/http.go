package usermetadatatransports

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	usermetadataendpoints "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/usermetadata/endpoints"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/state"
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

	getHandler := httptransport.NewServer(
		endpoints.GetUserMetadata,
		DecodeGetUserMetadataRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	return append(usermetadataServers, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: addHandler,
		Route:"/usermetadata",
		Methods: []string{"POST"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: getHandler,
		Route: "/usermetadata/{entityName}",
		Methods: []string{"GET"},
	})
}

func DecodeGetUserMetadataRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	if !state.GetState().IsAllowedToWrite() {
		return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
	}

	var req usermetadataendpoints.GetUserMetadataRequest

	decodedReq:= json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	if err != nil {
		return nil, err
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) <= 1 {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	} else if parts[len(parts) - 1] == "" {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	}

	req.EntityName = parts[len(parts) - 1]

	return req, nil
}