package mypackage

import (
	"errors"
	"fmt"
	"os"
)

// Logger 使用接口的方式实现一个既可以往终端写日志也可以往文件写日志的简易日志库
// 定义日志接口
type Logger interface {
	Log(msg string) error
}

// ConsoleLogger 实现两个版本的日志接口
// 控制台
type ConsoleLogger struct {
}

func (c *ConsoleLogger) Log(msg string) error {
	fmt.Println(msg)
	return nil
}

// FileLogger 文件
type FileLogger struct {
	filePath string
	file     *os.File
}

func NewFileLogger(filePath string) (*FileLogger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &FileLogger{filePath: filePath, file: file}, nil
}

func (f *FileLogger) Log(msg string) error {
	_, err := f.file.WriteString(msg + "\n")
	return err
}

// Close 关闭文件
func (f *FileLogger) Close() error {
	return f.file.Close()
}

// GetLogger 创建一个工厂函数，根据不同的需求返回不同类型的日志器。
func GetLogger(logType string, filePath string) (Logger, error) {
	switch logType {
	case "console":
		return &ConsoleLogger{}, nil
	case "file":
		if filePath == "" {
			return nil, errors.New("filePath cannot be empty for file logger")
		}
		return NewFileLogger(filePath)
	default:
		return nil, errors.New("unsupported log type")
	}
}
