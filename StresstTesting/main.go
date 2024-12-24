package main

import (
	"StressTesting/controller"
	"github.com/gin-gonic/gin"
)

//存储数据小案例

func main() {

	r := gin.Default()
	r.POST("/user", controller.ConnectMysql, controller.SaveUserHandler)
	r.Run(":8080")

}
