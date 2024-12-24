package routers

import (
	"GokitGRPCDemo/endpoints"
	"GokitGRPCDemo/middlewares"
	"GokitGRPCDemo/pb"
	"GokitGRPCDemo/transports"
	"context"
	GrpcTransport "github.com/go-kit/kit/transport/grpc"
	HttpTransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"net/http"
)

//利用Go-kit实现接口，Go-kit实现的接口，只用于数据透传，调用底层的endpoint，不做逻辑执行，最后再调用Service层代码

//数据透传层只是表面完成了接口，底层调用ServeGRPC接口。这个接口的底层实现就是decoder，endpoint，encoder。

// IGStringService 实现services.Service接口 也实现了proto定义的服务接口
// 既是Services接口类型，也是pb.StringServicesServer类型     注：应该这里只使用pb.StringServicesServer类型
type IGStringService struct {
	concat GrpcTransport.Handler
	diff   GrpcTransport.Handler
	pb.UnimplementedStringServicesServer
}

func (s IGStringService) Concat(ctx context.Context, request *pb.StringRequest) (*pb.StringResponse, error) {
	_, rep, err := s.concat.ServeGRPC(ctx, request)
	return rep.(*pb.StringResponse), err
}

func (s IGStringService) Diff(ctx context.Context, request *pb.StringRequest) (*pb.StringResponse, error) {
	_, rep, err := s.diff.ServeGRPC(ctx, request)
	return rep.(*pb.StringResponse), err
}

/*
	为transport层添加日志
*/

func NewRouter(endpoint endpoints.StringEndpoints) pb.StringServicesServer {
	return &IGStringService{
		concat: GrpcTransport.NewServer(endpoint.ConcatEndpoint, transports.DecodeString, transports.EncodeString),
		diff:   GrpcTransport.NewServer(endpoint.DiffEndpoint, transports.DecodeString, transports.EncodeString),
	}
}

//http

func NewHTTPServer(stringEndpoints endpoints.StringEndpoints, logger log.Logger) http.Handler {
	concatEndpoint := stringEndpoints.ConcatEndpoint
	//调用LoggerMiddleware函数返回一个函数，接着将concatEndpoint作为参数传递进去，返回值即为concatEndpoint的值。
	//返回的函数是Middleware类型的，Middleware函数类型是传递一个Endpoint类型的参数，返回一个Endpoint类型的参数。
	concatEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "concat"))(concatEndpoint)
	concatHandler := HttpTransport.NewServer(
		concatEndpoint,
		transports.DecodeHTTPString,
		transports.EncodeHTTPString,
	)
	diffEndpoint := stringEndpoints.DiffEndpoint
	diffEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "diff"))(diffEndpoint)
	diffHandler := HttpTransport.NewServer(
		diffEndpoint,
		transports.DecodeHTTPString,
		transports.EncodeHTTPString,
	)
	r := mux.NewRouter()
	r.Handle("/concat", concatHandler).Methods("POST")
	r.Handle("/diff", diffHandler).Methods("POST")
	return r
}
