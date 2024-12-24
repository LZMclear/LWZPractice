package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"grpcDemo2/pb"
	"net"
	"sync"
)

// grpcServer
type server struct {
	pb.UnimplementedWelcomeServer
	mu    sync.Mutex     //count的并发锁
	count map[string]int //记录每个name的请求次数
}

// 这个结构体主要提供名称解析与均衡负载中的服务
type serverA struct {
	pb.UnimplementedWelcomeServer
	Addr string
}

func (s *serverA) SayHelloA(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	reply := fmt.Sprintf("hello %s. [from %s]", in.GetName(), s.Addr)
	return &pb.HelloResponse{Reply: reply}, nil
}

var port = flag.Int("port", 8080, "服务端口")

// 名称解析与均衡负载主函数
func mainA() {
	flag.Parse()
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	//启动服务
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to listen, err:%v\n", err)
		return
	}
	//创建grpc服务
	s := grpc.NewServer()
	//注册服务
	pb.RegisterWelcomeServer(s, &serverA{Addr: addr})
	//启动服务
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve,err:%v\n", err)
		return
	}
}

// SayHello 这是proto定义需要实现的server服务的方法
func (s server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	s.mu.Lock()         //开局先上锁
	defer s.mu.Unlock() //函数运行结束后再解锁
	s.count[in.Name]++  //记录用户的请求次数

	//超过一次就返回错误
	if s.count[in.Name] > 1 {
		st := status.New(codes.ResourceExhausted, "request limit exceed") //这是一个status类型的错误
		//为错误添加详细信息
		details, err := st.WithDetails(&errdetails.QuotaFailure{Violations: []*errdetails.QuotaFailure_Violation{
			{
				Subject:     fmt.Sprintf("name:%s\n", in.Name),
				Description: "限制每个name只能被调用一次",
			},
		}})
		if err != nil {
			return nil, err
		}
		return nil, details.Err()
	}
	//正常返回响应
	reply := "hello" + in.GetName()
	return &pb.HelloResponse{Reply: reply}, nil
}

// SSLmain 自签名证书进行通信加密主函数
func SSLmain() {
	//启动监听服务
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("failed to listen")
		return
	}
	//加载server.crt证书和秘钥server.key
	cred, _ := credentials.NewServerTLSFromFile("./server.crt", "./server.key")
	//创建grpc服务
	s := grpc.NewServer(grpc.Creds(cred))
	//注册服务，同时初始化count
	pb.RegisterWelcomeServer(s, &server{count: make(map[string]int)})
	//启动服务
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve,err:%v\n", err)
		return
	}
}

func main() {
	mainA()
}
