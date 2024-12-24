package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// FprintPractice 将内容输出到一个io.Write类型的变量中通常用这个函数往文件中写入内容
func FprintPractice() {
	fmt.Fprint(os.Stdout, "像标准输出写入内容")
	fileObj, err := os.OpenFile("Fprint.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("打开文件出错，err:", err)
		return
	}
	name := "归途"
	// 向打开的文件句柄中写入内容
	fmt.Fprintf(fileObj, "往文件中写如信息：%s", name)
}

// SprintPractice 会把传入的数据生成一个字符串并返回
func SprintPractice() {
	s1 := fmt.Sprint("归途")
	name := "归途"
	age := 18
	s2 := fmt.Sprintf("name:%s,age:%d", name, age)
	s3 := fmt.Sprintln("归途")
	fmt.Println(s1, s2, s3)

}

// ErrofPractice 根据format参数生成格式化字符串并返回一个包含该字符串的错误
func ErrofPractice() {
	e := errors.New("原始错误e")
	w := fmt.Errorf("Wrap了一个错误%w", e)
	fmt.Println(w)
}

func bufioDemo() {
	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	fmt.Print("请输入内容：")
	text, _ := reader.ReadString('q') // 读到换行
	text = strings.TrimSpace(text)
	fmt.Printf("%#v\n", text)
}

func main() {
	now := time.Now()
	// 24小时制  最后两个展示的是星期和月份
	fmt.Println(now.Format("2006-01-02 15:04:05.000 Mon Jan"))
	// 12小时制
	fmt.Println(now.Format("2006-01-02 03:04:05.000 PM Mon Jan"))
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	timeObj, err := time.ParseInLocation("2006/01/02 15:04:05", "2022/10/05 11:25:20", loc)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(timeObj)
}
