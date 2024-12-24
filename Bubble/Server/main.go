package main

import (
	"Server/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//创建一个路由组
	v1Group := r.Group("/v1")
	//这个路由组全部调用连接数据库中间件
	v1Group.Use(controller.ConnectSql)
	//增添
	v1Group.POST("/todo", controller.HandlerPost)
	//删除
	v1Group.DELETE("/todo", controller.HandlerDelete)
	//修改
	v1Group.PUT("/todo", controller.HandlerUpdate)
	//查询
	v1Group.GET("/todo", controller.HandlerGet)
	r.Run(":9000")
}
