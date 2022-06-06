package util

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

func TransportLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time){
				logger.Log("transport_error", err, "took", time.Since(begin).Milliseconds())
			}(time.Now())

			return next(ctx, request)
		}
	}
}