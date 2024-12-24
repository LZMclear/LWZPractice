package main

import (
	"TrimService/ts"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strings"
)

var port = flag.Int("port", 8975, "service port")

type server struct {
	ts.UnimplementedTrimServer
}

// TrimSpace 去除字符串参数中的空格
func (s *server) TrimSpace(_ context.Context, req *ts.TrimRequest) (*ts.TrimResponse, error) {
	ov := req.GetS()
	v := strings.ReplaceAll(ov, " ", "")
	fmt.Printf("ov:%s v:%v\n", ov, v)
	return &ts.TrimResponse{S: v}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	ts.RegisterTrimServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
