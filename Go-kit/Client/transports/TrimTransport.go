package transports

import (
	"Client/endpoints"
	"Client/pb"
	"context"
)

// 请求与响应的转换函数

func EncodeTrimRequest(_ context.Context, response interface{}) (request interface{}, err error) {
	resp := response.(endpoints.TrimResponse)
	return &pb.TrimRequest{S: resp.S}, nil
}

func DecodeTrimResponse(_ context.Context, in interface{}) (interface{}, error) {
	resp := in.(*pb.TrimResponse)
	return endpoints.TrimResponse{S: resp.S}, nil
}
