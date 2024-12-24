package endpoints

import (
	"GokitRPCDemo/entity"
	"GokitRPCDemo/service"
	"context"
	"github.com/go-kit/kit/endpoint"
)

// 服务接口中只有一个方法，编写简单的适配器将服务中的每个方法转换为一个端点。
// 每个适配器接收一个AddServiceInter，并返回与其中一个方法对应的端点

func MakeSunEndpoints(svc service.AddServiceInter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(entity.SumRequest)
		sum, err := svc.Sum(ctx, req.A, req.B)
		if err != nil {
			return entity.SumResponse{V: sum, Err: err.Error()}, nil
		}
		return entity.SumResponse{V: sum}, nil
	}
}

func MakeConcatEndpoint(svc service.AddServiceInter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(entity.ConcatRequest)
		v, err := svc.Concat(ctx, req.A, req.B)
		if err != nil {
			return entity.ConcatResponse{V: v, Err: err.Error()}, nil
		}
		return entity.ConcatResponse{V: v}, nil
	}
}
