package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Args struct {
	X, Y int
}

type ServiceA struct {
}

// Add 为ServerA添加一个方法
func (s ServiceA) Add(arg *Args, reply *int) error {
	*reply = arg.X + arg.Y
	return nil
}

// 将ServerA类型注册为一个服务
func main() {
	service := new(ServiceA)
	rpc.Register(service) //注册RPC服务
	rpc.HandleHTTP()      //基于http协议
	l, e := net.Listen("tcp", ":9091")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Println("启动")
	http.Serve(l, nil)
}
