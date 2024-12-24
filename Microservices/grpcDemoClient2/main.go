package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"grpcDemoClient2/pb"
	"log"
	"time"
)

//grpcClient
//当服务端返回错误时，尝试从错误中获取详细信息

var name = flag.String("name", "gvto", "通过-name告诉server你是谁")

// ssl自签名客户端访问
func sslClient() {
	//解析命令行参数
	flag.Parse()
	//加载server.crt证书
	cred, _ := credentials.NewClientTLSFromFile("../grpcDemo2/server.crt", "")
	//连接server
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalf("grpc.Dail failed, err: %v\n", err)
		return
	}
	defer conn.Close()
	//创建客户端
	c := pb.NewWelcomeClient(conn)
	//开始调用rpc方法
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	response, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		//从错误中获取详细信息
		s := status.Convert(err) //将error类型转变为status类型
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *errdetails.QuotaFailure:
				fmt.Printf("Quota failure:%s\n", info)
			default:
				fmt.Printf("Unexpected type: %s\n", info)
			}
		}
		fmt.Printf("c.SayHello failed, err: %v\n", err)
		return
	}
	//没有错误，打印响应结果
	log.Printf("resp:%v\n", response.GetReply())
}

func main() {
	//初始化名称解析
	Init()
	LoadBalance()
}

// CustomClient 名称解析客户端
func CustomClient() {
	//连接server
	conn, err := grpc.Dial("gvto:///resolver.gvto.life", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc dial failed, err: %v\n", err)
		return
	}
	defer conn.Close() //最后关闭连接节省内存
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	//创建客户端
	c := pb.NewWelcomeClient(conn)
	//调用RPC方法
	response, err := c.SayHelloA(ctx, &pb.HelloRequest{Name: "lsl"})
	if err != nil {
		log.Printf("c.sayHelloA failed, err: %v\n", err)
		return
	}
	//输出
	fmt.Printf("response: %v\n", response.GetReply())
}

// LoadBalance 均衡负载调试
func LoadBalance() {
	//创建连接
	conn, err := grpc.Dial("gvto:///resolver.gvto.life",
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 这里设置初始策略
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	//创建客户端
	c := pb.NewWelcomeClient(conn)
	if err != nil {
		log.Fatalf("grpc dial failed, err: %v\n", err)
		return
	}
	//发起十次调用
	for i := 0; i < 10; i++ {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		response, err := c.SayHelloA(ctx, &pb.HelloRequest{Name: "lsl"})
		if err != nil {
			log.Printf("c.sayHello failed,err:%v\n", err)
			return
		}
		fmt.Printf("response: %v\n", response.GetReply())
		cancelFunc()
	}
}
