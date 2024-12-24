package main

import (
	"fmt"
	"mypackage"
)

func main() {
	exitChan := make(chan struct{ user string }, 1)
	users := struct {
		username string
	}{"gsy"}
	fmt.Println(users)
	fmt.Println(exitChan)
	exitChan <- struct {
		user string
	}{"lsl"}
LOOP:
	for {
		select {
		case <-exitChan:
			fmt.Println("运行的第一个")
			break LOOP
		default:
			fmt.Println("运行default")
		}
	}
	close(exitChan)
}

// 调用学生管理系统代码
func implementStudent() {
	//对学生进行初始化，初始化返回一个班级对象
	c := mypackage.Init()
	fmt.Println(c)
	//添加学生信息
	stu := &mypackage.Student{
		Id:   88,
		Name: "sss",
		Age:  25,
		Address: mypackage.Address{
			"河南",
			"新乡",
		},
	}
	c = c.AddStudent(stu)
	fmt.Println("添加学生后：")
	fmt.Println(c)
}

// 调用Logger代码
func implementLogger() {
	//获取控制台日志器
	consoleLogger, err := mypackage.GetLogger("console", "")
	if err != nil {
		fmt.Println("Error getting console logger:", err)
		return
	}

	// 获取文件日志器
	fileLogger, err := mypackage.GetLogger("file", "app.log")
	if err != nil {
		fmt.Println("Error getting file logger:", err)
		return
	}
	defer fileLogger.(*mypackage.FileLogger).Close()

	// 使用日志器
	consoleLogger.Log("This is a console log.")
	fileLogger.Log("This is a file log.")
}
