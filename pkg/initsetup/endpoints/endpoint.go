package initsetupendpoints

import (
	"context"

	initsetupservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/initsetup/service"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	AddConfig endpoint.Endpoint
	EditConfig endpoint.Endpoint
}

func New(svc initsetupservice.InitSetupService, logger log.Logger) Set {
	var addConfigEndpoint endpoint.Endpoint 
	{
		addConfigEndpoint = AddConfigEndpoint(svc)
		addConfigEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "addConfig"))(addConfigEndpoint)
	}

	var editConfigEndpoint endpoint.Endpoint
	{
		editConfigEndpoint = EditConfigEndpoint(svc)
		editConfigEndpoint = util.TransportLoggingMiddleware(log.With(logger, "method", "editConfig"))(editConfigEndpoint)
	}

	return Set{
		AddConfig: addConfigEndpoint,
		EditConfig: editConfigEndpoint,
	}
}

func AddConfigEndpoint(svc initsetupservice.InitSetupService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req:= request.(AddConfigRequest)

		config, err := svc.AddConfig(req.NewConfig)

		return AddConfigResponse{config, err}, err
	}
}

func EditConfigEndpoint(svc initsetupservice.InitSetupService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req:= request.(EditConfigRequest)

		isUpdated, err := svc.EditConfig(req.NewConfig)

		return EditConfigResponse{isUpdated, err}, err
	}
}