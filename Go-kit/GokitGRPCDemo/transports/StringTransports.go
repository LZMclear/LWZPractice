package transports

import (
	"GokitGRPCDemo/pb"
	"context"
	"encoding/json"
	"net/http"
)

//grpc的编码解码

func DecodeString(ctx context.Context, req interface{}) (request interface{}, err error) {
	request = req.(*pb.StringRequest)
	return
}

func EncodeString(ctx context.Context, rep interface{}) (response interface{}, err error) {
	response = rep.(*pb.StringResponse)
	return
}

//http的编码解码

func DecodeHTTPString(ctx context.Context, req *http.Request) (interface{}, error) {
	var request *pb.StringRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeHTTPString(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
