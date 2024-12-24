package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

// 用于签名的字符串
var mySingingKey = []byte("gvto")

// GenRegisteredClaims 使用默认声明claim创建jwt
func GenRegisteredClaims() (string, error) {
	//创建Claims
	Claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), //设置过期时间
		Issuer:    "Gvto",                                             //签发人
	}
	//生成Token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	//生成签名字符串
	return token.SignedString(mySingingKey)
}

// ValidateRegisteredClaims 解析jwt
func ValidateRegisteredClaims(tokenString string) bool {
	//解析Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySingingKey, nil
	})
	if err != nil { // 解析token失败
		return false
	}
	return token.Valid
}
func main() {
	r := gin.Default()
	r.POST("/authority", authHandler)
	//验证
	r.GET("/check", JWTAuthMiddleware(), checkHandler)
	r.Run(":8080")

}

func checkHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"username": username},
	})
}
