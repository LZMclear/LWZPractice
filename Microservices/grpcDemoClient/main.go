package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"grpcDemoClient/pb"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// hello client

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "127.0.0.1:8080", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func runLotsOfReplies(c pb.GreeterClient) {
	//server端流式rpc
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	replies, err := c.LotsOfReplies(ctx, &pb.HelloRequest{
		Name: defaultName,
	})
	if err != nil {
		log.Fatalf("c.LotsOfReplies failed, err: %v\n", err)
		return
	}
	for {
		//接收服务端返回的流式数据收到io.EOF或者错误时退出
		recv, err := replies.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("c.LotsOfReplies failed, err: %v\n", err)
			return
		}
		log.Printf("got replies: %q\n", recv.GetReply())
	}
}

func runLotsOfGreetings(c pb.GreeterClient) {
	//创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.LotsOfGreetings(ctx)
	if err != nil {
		log.Fatalf("c.LotsGreetings failed, err: %v\n", err)
		return
	}
	names := []string{
		"lsl", "gsy", "zyl",
	}
	for _, name := range names {
		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("c.LotsOfGreetings stream.Send(%v) failed, err: %v", name, err)
		}
	}
	recv, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("c.LotsOfGreetings failed: %v", err)
	}
	log.Printf("got reply: %v", recv.GetReply())
}

// 简单的对话方式，先获取服务器返回的字符串，将
func runBindHello(c pb.GreeterClient) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelFunc()
	//双向流模式
	stream, err := c.BindHello(ctx)
	if err != nil {
		log.Fatalf("c.BindHello failed! err: %v\n", err)
		return
	}
	waitc := make(chan struct{})
	go func() { //开了一个线程接收服务器的响应并打印出来
		for {
			//接收服务端返回的响应
			recv, err2 := stream.Recv()
			if err2 == io.EOF {
				close(waitc)
				return
			}
			if err2 != nil {
				log.Fatalf("c.BindHello stream.Recv failed, err: %v\n", err)
			}
			fmt.Printf("AI: %s\n", recv.GetReply())
		}
	}()
	//从标准输入获取用户输入
	reader := bufio.NewReader(os.Stdin)
	for {
		readString, _ := reader.ReadString('\n') //读到换行
		readString = strings.TrimSpace(readString)
		if len(readString) == 0 {
			continue
		} //字符串长度为零，跳过下面代码片段，继续读取
		if strings.ToUpper(readString) == "QUIT" {
			break
		}
		//将获取到的数据发送服务器
		if err := stream.Send(&pb.HelloRequest{Name: readString}); err != nil {
			log.Fatalf("c.BindHello steam.Send(%v) failed, err:%v\n", readString, err)
		}
	}
	stream.CloseSend()
	<-waitc
}

// 普通RPC调用客户端metadata操作
func unaryCallMetadata(c pb.GreeterClient, name string) {
	//创建metadata
	md := metadata.Pairs(
		"token", "app-test-gvto",
		"request_id", "1234567890",
	)
	newmd := metadata.New(map[string]string{"sss": "sss", "aaa": "bbb"})
	//基于metadata创建ctx
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Join(md, newmd))
	//RPC调用
	var header, trailer metadata.MD
	response, err := c.SayHello(ctx, &pb.HelloRequest{Name: name}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatalf("c.SayHello failed, err: %v\n", err)
		return
	}
	//从header中提取location
	if v, ok := header["location"]; ok {
		fmt.Printf("location from header:\n")
		for i, e := range v {
			fmt.Printf("%d, %s\n", i, e)
		}
	} else {
		log.Printf("location expected but doesn't exist in header")
	}
	//获取响应结果
	fmt.Printf("get response: %s\n", response.Reply)
	//从trailer中提取timestamp(好像是时间戳)
	if v, ok := trailer["timestamp"]; ok {
		fmt.Printf("timestamp from trailer:\n")
		for i, e := range v {
			fmt.Printf("%d, %s\n", i, e)
		}
	} else {
		log.Printf("timestamp expected but doesn't exist in trailer")
	}
}

// StreamCallMetadata 流式RPC调用客户端metadata操作
func StreamCallMetadata(c pb.GreeterClient, name string) {
	//创建带有metadata的context
	md := metadata.Pairs("token", "app-test-gvto")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	//使用带有metadata的context执行RPC调用 token数据存储在ctx中调用服务端的函数，返回一个数据流
	stream, err := c.BindHelloWithMetadata(ctx)
	if err != nil {
		log.Fatalf("执行c.BindHelloWithMetadata 失败, err: %v\n", err)
		return
	}
	//开一个线程接收服务器的响应
	go func() {
		//当header到达时读取header
		header, err2 := stream.Header()
		if err2 != nil {
			log.Fatalf("failed to get header from stream: %v", err)
		}
		//从返回响应的header中读取数据
		if l, ok := header["location"]; ok {
			fmt.Printf("location from header:\n")
			for i, e := range l {
				fmt.Printf(" %d. %s\n", i, e)
			}
		} else {
			log.Println("location expected but doesn't exist in header")
			return
		}
		//发送所有的请求到server
		for i := 0; i < 5; i++ {
			if err := stream.Send(&pb.HelloRequest{Name: name}); err != nil {
				log.Fatalf("failed to send streaming: %v\n", err)
			}
		}
		//发送完关闭数据流
		stream.CloseSend()
	}()
	//读取所有的响应
	var rpcStatus error
	fmt.Println("get response:\n")
	for {
		recv, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		fmt.Printf(" - %s\n", recv.Reply)
	}
	if rpcStatus != io.EOF { //错误不是终止符号输出具体错误
		log.Printf("failed to finish server streaming: %v", rpcStatus)
		return
	}
	// 当RPC结束时读取trailer
	trailer := stream.Trailer()
	// 从返回响应的trailer中读取metadata.
	if t, ok := trailer["timestamp"]; ok {
		fmt.Printf("timestamp from trailer:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Printf("timestamp expected but doesn't exist in trailer")
	}
}

func main() {
	flag.Parse()
	//连接到Server端，此处禁用安全传输
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failed, err: %v\n", err)
		return
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	StreamCallMetadata(c, "gvto")
}
