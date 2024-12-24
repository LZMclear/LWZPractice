package main

import (
	"GokitRPCDemo/dao"
	"GokitRPCDemo/endpoints"
	"GokitRPCDemo/entity"
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"log"
	"net/http"
)

func decodeSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request entity.SumRequest
	//json.NewDecoder创建一个新的json解码器，r.Body是实现了io.Reader接口的HTTP请求体，表示接收的原始数据
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request entity.ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// 将返回的数据编码为json格式
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {
	svc := dao.AddService{} //结构体实例化  可以用AddServiceInter类型的接口表示
	sumHandler := httptransport.NewServer(
		endpoints.MakeSunEndpoints(svc),
		decodeSumRequest,
		encodeResponse,
	)

	concatHandler := httptransport.NewServer(
		endpoints.MakeConcatEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)
	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
