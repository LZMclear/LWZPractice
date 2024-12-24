package middlewares

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

/*
	endpoint.Middleware中middleware类型是一个接受Endpoint返回Endpoint的函数
	1. 第一层返回Middleware类型的函数
	2. 第二层返回Endpoint类型的函数
*/

func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", "calling endpoint")
			defer logger.Log("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}
