package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 定义一个中间件，中间件就是一个handlerFunc类型的函数，传递一个*gin.Context类型的参数的函数就是handlerFunc函数
func middleFunc(c *gin.Context) {
	//计时
	start := time.Now()
	//执行后续的处理函数，也就是middleFunc后面的函数
	c.Next()
	cost := time.Since(start)
	fmt.Printf("执行处理函数花费时间：%v\n", cost)
}

func main() {
	r := gin.Default()
	//如果有多个路由使用中间件，可以使用r.Use(...)  全局使用中间件
	r.GET("/index", middleFunc, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.Run(":8080")
}
