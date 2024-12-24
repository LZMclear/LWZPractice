package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

var logger *zap.Logger

// LogConfig 定义日志相关配置
type LogConfig struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"maxsize"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}

// Config 整个项目的配置
type Config struct {
	Mode       string `json:"mode"`
	Port       int
	*LogConfig `json:"log"`
}

// Conf 定义全局配置变量
var Conf = new(Config)

// 在项目中先从配置文件中加载配置信息，再调用InitLogger完成logger实例的初始化
// 通过r.Use(logger.GinLogger(), logger.GinRecovery(true))注册我们的中间件来使用zap接收gin框架自身的日志
func main() {
	//从config.json中加载配置
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("读取配置文件出错")
		return
	}
	json.Unmarshal(b, Conf)
	if err := InitLogger(Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err %v\n", err)
		return
	}
	gin.SetMode(Conf.Mode)

	r := gin.New()
	r.Use(GinLogger(), GinRecovery(true))
	r.GET("/hello", func(c *gin.Context) {
		var name = "lsl"
		var age = 18
		logger.Debug("this is hello func", zap.String("user", name), zap.Int("age", age))
		c.String(http.StatusOK, "hello %s", "lsl")
	})
	address := fmt.Sprintf(":%v", Conf.Port)
	r.Run(address)
}
