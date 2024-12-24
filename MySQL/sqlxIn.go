package main

import (
	"database/sql/driver"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Value 使用sqlx.In批量插入数据，首先使结构体实现driver.Valuer()接口
func (user test) Value() (driver.Value, error) {
	//返回一个interface类型的切片，切片中存储的是user的值
	return []interface{}{user.Username, user.Password}, nil
}

// BatchInsertUser 使用sqlxIn插入数据  sqlx.In帮我们拼接语句和参数
func BatchInsertUser(users []interface{}) {
	query, args, _ := sqlx.In("INSERT into test (username, password) values (?),(?),(?)", users...)
	fmt.Println(query)
	fmt.Println(args)
	_, err := sqlxdb.Exec(query, args...)
	if err != nil {
		fmt.Printf("批量插入数据失败，错误：%v\n", err)
		return
	}
}

// BatchInsertUsers3 values后面匹配的字段名需要配合结构体中为各个变量配置tag使用
func BatchInsertUsers3(users []*test) {
	_, err := sqlxdb.NamedExec("INSERT INTO test (username, password) VALUES (:Username, :Password)", users)
	if err != nil {
		fmt.Printf("批量插入数据失败，错误：%v\n", err)
		return
	}
}
