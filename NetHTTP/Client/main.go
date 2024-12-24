package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	PostSample()
}

func PostSample() {
	url := "http://127.0.0.1:9090/post"
	// 表单数据
	contentType := "application/x-www-form-urlencoded"
	data := "name=小王子&age=18"
	// json
	//contentType := "application/json"
	//data := `{"name":"小王子","age":18}`
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}
	fmt.Println(string(b))
}

func GetSample() {
	apiUrl := "http://127.0.0.1:9090/get"
	// URL 参数
	data := url.Values{}
	data.Set("name", "小王子")
	data.Set("age", "18")
	//解析原始URL字符串
	u, err := url.Parse(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
	}
	u.RawQuery = data.Encode() // URL encode
	fmt.Println(u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}
	fmt.Println(string(b))
}
