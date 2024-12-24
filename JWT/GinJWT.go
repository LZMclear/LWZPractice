package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type User struct {
	Username string `json:"user"`
	Password string `json:"pwd"`
}

func authHandler(c *gin.Context) {
	//用户发送用户名和密码
	var user User
	//使用ShouldBind绑定参数
	err := c.ShouldBind(&user)
	fmt.Println(c.PostForm("password"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "无效的参数",
		})
	}
	fmt.Printf("绑定的user信息：%v\n", user)
	//校验用户名和密码是否正确
	if user.Username == "gvto" && user.Password == "251210" {
		//生成Token
		tokenString, _ := GenToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "success",
			"data": tokenString,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "鉴权失败",
	})
}

// JWTAuthMiddleware 实现一个校验token的中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//客户端携带Token有三种方式，1.放在请求头 2.放在请求体 3.放在URI
		//这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		fmt.Printf("authHeader:%v\n", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		//按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := parseToken(parts[1])
		fmt.Printf("mc%v\n", mc)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.username)
		fmt.Println(mc.username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
