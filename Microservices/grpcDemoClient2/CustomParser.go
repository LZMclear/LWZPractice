package main

import "google.golang.org/grpc/resolver"

//自定义name resolver

const (
	myScheme   = "gvto"
	myEndPoint = "resolver.gvto.life"
)

var addr = []string{"127.0.0.1:8972", "127.0.0.1:8973"}

// gvtoResolver  自定义name resolver，实现Resolver接口
type gvtoResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string //这里面应该存储的是IP地址
}

func (r *gvtoResolver) ResolveNow(o resolver.ResolveNowOptions) {
	//addrsStore是map类型，[]里面是字符串，取出键为r.target.Endpoint()的值，值类型为一个字符串切片，里面装的是IP地址
	//addrStrs为存储了IP地址的字符串切片
	addrStrs := r.addrsStore[r.target.Endpoint()]
	//初始化一个存储resolver.Address类型的切片，长度大小为上面字符串切片的长度
	addrList := make([]resolver.Address, len(addrStrs))
	//随后将上面字符串类型的IP地址转换为解析器地址类型存储到addrList切片中
	for i, str := range addrStrs {
		addrList[i] = resolver.Address{Addr: str}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*gvtoResolver) Close() {}

// gvtoResolverBuilder 需要实现Builder接口
// 也就是Builder接口中有一个Bind方法 需要gvtoResolverBuilder结构体实现这个方法 也就是实现了这个接口
// 所以gvtoResolverBuilder类型可以用Builder接口表示
type gvtoResolverBuilder struct{}

func (*gvtoResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opt resolver.BuildOptions) (resolver.Resolver, error) {
	r := &gvtoResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			myEndPoint: addr,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (*gvtoResolverBuilder) Scheme() string {
	return myScheme
}

func Init() {
	//注册gvtoResolverBuilder
	//gvtoResolverBuilder实现了Build方法，就是builder接口类型
	resolver.Register(&gvtoResolverBuilder{})
}

/*
	运行流程:
		grpc.Dail会根据gvto这个scheme找到我们注册的gvtoResolverBuilder，调用bind方法构建自定义gvtoResolver
*/
