package main

import (
	"GokitGRPCDemo/endpoints"
	"GokitGRPCDemo/pb"
	"GokitGRPCDemo/routers"
	"GokitGRPCDemo/services"
	"context"
	"flag"
	"fmt"
	goLog "github.com/go-kit/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	httpAddr = flag.String("http-addr", ":8972", "HTTP listen address")
	grpcAddr = flag.String("grpc-addr", ":8080", "gRPC listen address")
)

func main() {
	//初始化logger
	logger := goLog.NewLogfmtLogger(os.Stderr)
	var g errgroup.Group
	svc := services.IStringService{Logger: logger}
	Endpoints := endpoints.StringEndpoints{
		ConcatEndpoint: endpoints.MakeConcatEndpoint(svc),
		DiffEndpoint:   endpoints.MakeDiffEndpoint(svc),
	}
	g.Go(func() error {
		//grpc服务
		listen, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			log.Fatalf("listen port failed, err: %v\n", err)
			return err
		}
		defer listen.Close()
		grpcServer := grpc.NewServer()
		r := routers.NewRouter(Endpoints)
		pb.RegisterStringServicesServer(grpcServer, r)
		return grpcServer.Serve(listen)
	})
	g.Go(func() error {
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			fmt.Printf("http: net.Listen(tcp, %s) failed, err:%v\n", *httpAddr, err)
			return err
		}
		defer httpListener.Close()
		httpHandler := routers.NewHTTPServer(Endpoints, logger)
		return http.Serve(httpListener, httpHandler)
	})
	//创建一个客户端调用
	g.Go(func() error {
		conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("connect failed, err: %v\n", err)
			return err
		}
		defer conn.Close()
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()
		servicesClient := pb.NewStringServicesClient(conn)
		response, err := servicesClient.Diff(ctx, &pb.StringRequest{
			A: "B",
			B: "B",
		})
		fmt.Println(response)
		return err
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("server exit with err:%v\n", err)
	}

}
