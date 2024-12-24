package endpoints

import (
	"Client/pb"
	"Client/transports"
	"github.com/go-kit/kit/endpoint"
	GrpcTransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type TrimRequest struct {
	S string
}

type TrimResponse struct {
	S string
}

// 这个端点与GokitGRPCDemo中的端点定义不同，这个是客户端的endpoint
func makeTrimEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return GrpcTransport.NewClient(
		cc,
		"ts.Trim",
		"TrimSpace",
		transports.EncodeTrimRequest,
		transports.DecodeTrimResponse,
		pb.TrimResponse{},
	).Endpoint()
}
