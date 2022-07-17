package usermetadatatransports

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
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

	getHandler := httptransport.NewServer(
		endpoints.GetUserMetadata,
		DecodeGetUserMetadataRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	getMetadataByEntityHandler := httptransport.NewServer(
		endpoints.GetMetadataByEntity,
		util.DecodeGenericHTTPIdentifierRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	getMetadataByEntityPagedHandler := httptransport.NewServer(
		endpoints.GetMetadataByEntityPaged,
		DecodeGetMetadataByEntityPagedRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	updateMetadadaHandler := httptransport.NewServer(
		endpoints.UpdateMetadata,
		DecodeUpdateMetadataRequest,
		util.EncodeHTTPGenericResponse,
		serverOptions...
	)

	deleteMetadataHandler := httptransport.NewServer(
		endpoints.DeleteMetadata,
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
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: getMetadataByEntityHandler,
		Route: "/usermetadata/all/{entityName}",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: getMetadataByEntityPagedHandler,
		Route: "/usermetadata/paged/{entityName}",
		Methods: []string{"GET"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: updateMetadadaHandler,
		Route: "/usermetadata/{entityName}",
		Methods: []string{"PUT"},
	}, commonmodels.HttpServerConfig{
		NeedsAuth: true,
		Server: deleteMetadataHandler,
		Route: "/usermetadata/{entityName}",
		Methods: []string{"DELETE"},
	})
}

func DecodeGetUserMetadataRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req usermetadataendpoints.GetUserMetadataRequest

	decodedReq:= json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req)

	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) <= 1 {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	} else if parts[len(parts) - 1] == "" {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	}

	req.EntityName = parts[len(parts) - 1]

	if err == io.EOF {
		req.Query = map[string]interface{}{}
		return req, nil
	}

	if err != nil {
		return nil, err
	}

	return req, nil
}

func DecodeUpdateMetadataRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	requestScope := ctx.Value("requestScope").(commonmodels.RequestScope)

	if !requestScope.IsAllowedToUpdate() {
		return nil, errors.New(constants.UNAUTHORIZED_ATTEMPT)
	}
	
	var req usermetadataendpoints.UpdateMetadataRequest

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

func DecodeGetMetadataByEntityPagedRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req usermetadataendpoints.GetMetadataByEntityPagedRequest

	queries := r.URL.Query()

	if queries.Get("page") != ""{
		parsedPage, parseErr := strconv.Atoi(queries.Get("page"))
		if parseErr != nil {
			req.Page = 1
		}else{
			req.Page = parsedPage
		}
	} else{
		req.Page = 1
	}

	if queries.Get("size") != ""{
		parsedSize, parseErr := strconv.Atoi(queries.Get("size"))
		if parseErr != nil {
			req.Size = 5
		}else{
			req.Size = parsedSize
		}
	} else{
		req.Size = 5
	}
	
	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) <= 1 {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	} else if parts[len(parts) - 1] == "" {
		return nil, errors.New(constants.RECORD_NOT_FOUND)
	}

	req.EntityName = parts[len(parts) - 1]

	decodedReq := json.NewDecoder(r.Body)

	decodedReq.DisallowUnknownFields()

	err := decodedReq.Decode(&req.Query)

	if err == io.EOF {
		req.Query = map[string]interface{}{}
		return req, nil
	}

	if err != nil {
		return req, err
	}

	return req, nil
	
}