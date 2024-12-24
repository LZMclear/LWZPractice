package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	SetupLogger()
	sampleHttp("www.baidu.com")
	sampleHttp("http://www.baidu.com")
}

// SetupLogger 实现Go Logger
func SetupLogger() {
	logFiler, _ := os.OpenFile("./test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	log.SetOutput(logFiler)
}

// 建立一个到URL的HTTP连接，并将状态代码/错误记录到日志文件中
func sampleHttp(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Erro fetching url %s: %s", url, err.Error())
	} else {
		log.Printf("Status Code for %s : %s", url, resp.Status)
		resp.Body.Close()
	}
}
