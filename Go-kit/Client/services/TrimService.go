package services

import (
	"Client/pb"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type withTrimMiddleware struct {
	TrimService endpoint.Endpoint
}

// Sum 提供具体的实现
func (mw withTrimMiddleware) Sum(ctx context.Context, in *pb.StrRequest) (*pb.StrResponse, error) {
	return &pb.StrResponse{R: in.A + in.B}, nil
}

// Concat 发起grpc调用外部trim service服务
func (mw withTrimMiddleware) Concat(ctx context.Context, in *pb.StrRequest) (*pb.StrResponse, error) {
	//先调用trim服务，去除字符串中存在的空格

	return nil, nil
}
