package entity

// SumRequest Sum方法的参数.
type SumRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

// SumResponse Sum方法的响应
type SumResponse struct {
	V   int    `json:"v"`
	Err string `json:"err,omitempty"`
}

// ConcatRequest Concat方法的参数.
type ConcatRequest struct {
	A string `json:"a"`
	B string `json:"b"`
}

// ConcatResponse  Concat方法的响应.
type ConcatResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}
