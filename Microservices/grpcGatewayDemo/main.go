package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpcGatewayDemo/pb"
	"log"
	"net"
	"net/http"
	"strings"
)

type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello 实现SayHello方法
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("传入的参数为：%s\n", in.GetName())
	return &pb.HelloResponse{Reply: "message:" + in.Name}, nil
}

func main() {
	Demo()
}

/*
访问http服务，gwServer将请求交给gwmu处理
grpc的服务处理程序通过RegisterGreeterHandler已经注册到了gwmu中
gwmu可以处理http请求转发到tcp服务器中
*/

func Demo() {
	//监听tcp连接
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("net.Listen failed, err: %v\n", err)
		return
	}
	//创建服务
	s := grpc.NewServer()
	//注册服务
	pb.RegisterGreeterServer(s, &server{})
	//启动服务
	log.Println("Service grpc on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(listen))
	}()
	//使用线程挂起服务，开启服务后直接创建客户端进行连接，没有另开一个客户端
	//上面内容创建正常的tcp服务

	//创建一个客户端连接
	//grpc-gateway就是通过这个代理请求（将http请求转换为tcp请求）
	conn, err := grpc.DialContext(context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(), //确保在连接成功之前阻塞调用，直到建立连接
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux() //创建一个新的grpcGateway多路复用器Mux，用于处理http请求并将其转发到tcp服务器上
	//将grpc服务的处理程序注册到grpcGateway多路复用器
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	//创建http服务器，监听8090端口，将请求交给gwmux处理
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
	// 8090端口提供gRPC-Gateway服务
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe()) //启动 HTTP 服务器，开始监听请求并处理它们。
}

// Demo2 同一端口提供http API和tcp API
func Demo2() {
	//创建一个tcp监听端口
	listen, err := net.Listen("tcp", ":8091")
	if err != nil {
		log.Fatalf("listen tcp port failed, err: %v\n", err)
		return
	}

	//创建grpc服务对象
	s := grpc.NewServer()
	//将自定义的greeter服务注册到grpc服务器上 这个greeter服务就是实例化的server
	pb.RegisterGreeterServer(s, &server{})

	//grpc gateway mux
	gwmux := runtime.NewServeMux()                                                      //创建一个http Mux 用于处理grpc-gateway请求
	dops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())} //指定使用不安全的连接
	//将grpc服务的http处理程序注册到Mux，以便通过 HTTP/1.1 调用 gRPC 服务
	err = pb.RegisterGreeterHandlerFromEndpoint(context.Background(), gwmux, "127.0.0.1:8091", dops)
	if err != nil {
		log.Fatalf("failed to register gwmux, err: %v\n", err)
		return
	}

	//创建一个新的 HTTP Mux，并将 gRPC-Gateway 的 Mux 注册到根路径 /。
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	//定义http server配置
	gwServer := &http.Server{
		Addr:    "127.0.0.1:8091",
		Handler: grpcHandlerFunc(s, mux),
	}
	log.Println("Serving on http://127.0.0.1:8091")
	log.Fatalln(gwServer.Serve(listen)) // 启动HTTP服务
}

// grpc请求和http请求分别调用不同的handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.ProtoMajor == 2 && strings.Contains(request.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(writer, request)
		} else {
			otherHandler.ServeHTTP(writer, request)
		}
	}), &http2.Server{})
}
