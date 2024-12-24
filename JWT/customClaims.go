package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// CustomsClaims 自定义声明类型，并内嵌默认的Claims
type CustomsClaims struct {
	username             string `json:"username"`
	jwt.RegisteredClaims        //内嵌标准的声明
}

// TokenExpireTime 定义jwt过期的时间
const TokenExpireTime = time.Hour * 24

// CustomSecret 定义一个用于签名的字符串
var CustomSecret = []byte("gsy")

// GenToken 生成jwt
func GenToken(username string) (string, error) {
	//创建自定义声明
	claims := CustomsClaims{
		username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireTime)),
			Issuer:    "lsl",
		},
	}
	//使用指定的签名方法创建签发对象  将claims作为jwt的负载
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//使用指定的秘钥（CustomSecret）签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// 解析jwt
func parseToken(tokenString string) (*CustomsClaims, error) {
	//解析token
	//如果是自定义Claim结构体则需要使用ParseWithClaims方法   &CustomsClaims{}加上&是初始化结构体
	token, err := jwt.ParseWithClaims(tokenString, &CustomsClaims{}, func(token *jwt.Token) (interface{}, error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	//对token对象中的claim进行类型断言
	if claims, ok := token.Claims.(*CustomsClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
