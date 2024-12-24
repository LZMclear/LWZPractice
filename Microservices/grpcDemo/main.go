package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpcDemo/pb"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

//server端

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	words := []string{
		"你好",
		"Hello",
		"こんにちは",
		"안녕하세요",
	}
	for _, word := range words {
		data := &pb.HelloResponse{
			Reply: word + in.GetName(),
		}
		//使用send方法返回多个数据
		if err := stream.Send(data); err != nil {
			return err
		}

	}
	return nil
}

// LotsOfGreetings 服务端接收流数据
func (s *server) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
	reply := "你好"
	//接收客户端的流式数据
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			//说明接收完毕，统一回复
			return stream.SendAndClose(&pb.HelloResponse{
				Reply: reply,
			})
		}
		if err != nil {
			return err
		}
		reply += recv.GetName()
	}
}

// BindHello 双向流数据服务
func (s *server) BindHello(stream pb.Greeter_BindHelloServer) error {
	for {
		//接收流式请求
		recv, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		reply := magic(recv.GetName()) //对收到的数据进行处理
		//返回流式响应
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

// 处理数据的magic函数
func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

// SayHello 普通RPC调用服务端metadata操作
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	//通过defer设置trailer
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		grpc.SetTrailer(ctx, trailer)
	}()
	//从客户端请求上下文中读取metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "UnarySayHello failed to get metadata")
	}
	if t, ok := md["token"]; ok {
		fmt.Printf("token from metadata:\n")
		if len(t) < 1 || t[0] != "app-test-gvto" {
			return nil, status.Error(codes.Unauthenticated, "认证失败")
		}
	}
	//创建和发送header
	header := metadata.New(map[string]string{"location": "Beijing"})
	grpc.SendHeader(ctx, header)
	fmt.Printf("request received: %v, say hello...\n", in)
	return &pb.HelloResponse{Reply: in.Name}, nil
}

// BindHelloWithMetadata 流式RPC调用metadata操作
func (s *server) BindHelloWithMetadata(stream pb.Greeter_BindHelloWithMetadataServer) error {
	//在defer中创建trailer记录返回时间
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		stream.SetTrailer(trailer)
	}()
	//从client中读取metadata
	md, b := metadata.FromIncomingContext(stream.Context())
	if !b {
		return status.Errorf(codes.DataLoss, "BidirectionalStreamingSayHello: failed to get metadata")
	}
	//从metadata中读取token
	if t, ok := md["token"]; ok {
		fmt.Printf("token from metadata:\n")
		//以token为键的值可能有多个，遍历输出
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}
	//创建和发送header
	header := metadata.New(map[string]string{"location": "X2Q"})
	stream.SendHeader(header)
	//读取请求数据发送响应数据
	for {
		recv, err := stream.Recv()
		if err != nil {
			fmt.Printf("get request failed, err：%v\n", err)
		}
		fmt.Printf("request received %v\n sending reply", recv)
		if err := stream.Send(&pb.HelloResponse{Reply: recv.Name}); err != nil {
			return err
		}
	}
}
func main() {
	//监听本地8080端口
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("listen failed, err: %v\n", err)
		return
	}
	s := grpc.NewServer()                  //创建grpc服务器
	pb.RegisterGreeterServer(s, &server{}) //在grpc服务端注册服务
	//启动服务
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("failed to server: %v", err)
		return
	}

}
