package endpoints

import (
	"GokitGRPCDemo/pb"
	"GokitGRPCDemo/services"
	"context"
	"github.com/go-kit/kit/endpoint"
)

// 服务接口中只有一个方法，编写简单的适配器将服务中的每个方法转换为一个端点。

type StringEndpoints struct {
	ConcatEndpoint endpoint.Endpoint
	DiffEndpoint   endpoint.Endpoint
}

//传递的参数是接口类型的，这样可以使所有实现这个接口的结构体都可以使用这个端点

func MakeConcatEndpoint(svc services.StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.StringRequest)
		response, err = svc.Concat(ctx, req)
		return
	}
}

func MakeDiffEndpoint(svc services.StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.StringRequest)
		response, err = svc.Diff(ctx, req)
		return
	}
}
