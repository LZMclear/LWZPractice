package services

import (
	"GokitGRPCDemo/pb"
	"context"
	"errors"
	"github.com/go-kit/log"
	"time"
)

//实现底层service各个接口的逻辑

// StringService 创建字符串服务接口
type StringService interface {
	Concat(ctx context.Context, request *pb.StringRequest) (*pb.StringResponse, error)
	Diff(ctx context.Context, request *pb.StringRequest) (*pb.StringResponse, error)
}

// IStringService 结构体实现接口 里面添加日志字段，用来打印详细请求参数
type IStringService struct {
	Logger log.Logger
}

func (s IStringService) Concat(_ context.Context, request *pb.StringRequest) (*pb.StringResponse, error) {
	//添加日志功能
	defer func(begin time.Time) {
		s.Logger.Log("method", "concat", "A", request.A, "B", request.B, "took", time.Since(begin))
	}(time.Now())
	if len(request.A+request.B) > 10 {
		return nil, errors.New("too long strings")
	}
	return &pb.StringResponse{Msg: request.A + request.B}, nil
}

func (s IStringService) Diff(_ context.Context, request *pb.StringRequest) (*pb.StringResponse, error) {
	//添加日志功能
	defer func(begin time.Time) {
		s.Logger.Log("method", "diff", "A", request.A, "B", request.B, "took", time.Since(begin))
	}(time.Now())
	if len(request.A+request.B) > 10 {
		return nil, errors.New("too long strings")
	}
	if request.A == request.B {
		return &pb.StringResponse{Msg: "two strings is same"}, nil
	} else {
		return &pb.StringResponse{Msg: "two strings is not same"}, nil
	}
}
