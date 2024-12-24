package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func main() {
	customizedLogger()
	defer logger.Sync()
	for i := 0; i < 3000; i++ {
		SugarHttpGet("www.baidu.com")
		SugarHttpGet("http://www.baidu.com")
	}
}

// InitLogger 通过调用NewProduction，NewDevelopment，Example创建一个Logger，三种方法的记录的信息不同
func InitLogger() {
	logger, _ = zap.NewProduction()
}

func InitSugarLogger() {
	logger, _ = zap.NewProduction()
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Error fetch url...", zap.String("url", url), zap.Error(err))
	} else {
		logger.Info("Success...", zap.String("StatusCode", resp.Status), zap.String("url", url))
		resp.Body.Close()
	}
}

func SugarHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s\n", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s\n", resp.Status, url)
		resp.Body.Close()
	}
}

// 定制logger
func customizedLogger() {
	//使用zap.New()方法传递所有配置
	//将文件写到哪，对文件的一些操作
	writeSyncer := getLogWriter()
	//如何写入日志，日志的格式编辑
	encoder := getEncoder()
	//Log Level:那种级别的日志被写入
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	//将调用函数信息添加到日志
	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

// 将文件写到哪
func getLogWriter() zapcore.WriteSyncer {
	//利用io.MultiWriter支持文件和终端两个输出目标
	//ws := io.MultiWriter(file, os.Stdout)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log", //日志文件的位置
		MaxSize:    1,            //进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,            //保留旧文件的最大个数
		MaxAge:     30,           //保留旧文件的最大天数
		Compress:   false,        //是否压缩归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 如何写入日志
func getEncoder() zapcore.Encoder {
	//这个应该是以JSON格式的方式被写入
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	//修改时间编码器，在日志文件中使用大写字母记录日志级别
	//覆盖默认的ProductionEncoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
