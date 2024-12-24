package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化Logger
func InitLogger(cf *LogConfig) (err error) {
	writeSyncer := getLogWriter(cf.Filename, cf.MaxSize, cf.MaxBackups, cf.MaxAge)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cf.Level)) //控制的是以哪种级别的日志被写入
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)
	logger = zap.New(core, zap.AddCaller())
	//zap.ReplaceGlobals(logger) //替换zap包中全局的logger实例，后续在其他包中只需调用zap.L即可
	return
}

// 日志的写入方式  对文件内容的操作
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig) //NewJSON...是以json格式写入
}

// 日志写哪，以及各种配置  对文件的操作
func getLogWriter(filename string, MaxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	luberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    MaxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(luberJackLogger)
}
