package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Email    *string //指向字符串的指针
}

func main() {
	//连接数据库
	dsn := "root:251210@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	sqlDB, _ := sql.Open("mysql", dsn) //创建数据库连接
	//通过已有的数据库连接来初始化gorm.DB
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}))

	//
	sql, err := db.DB()
	//设置空闲连接池的最大连接数
	sql.SetMaxIdleConns(10)
	//设置数据库的最大打开连接数
	sql.SetMaxOpenConns(100)
	//设置连接可以重用的最大时间量
	sql.SetConnMaxLifetime(time.Hour)
	if err != nil {
		fmt.Printf("连接数据库错误:%v\n", err)
		return
	}
	//自动迁移模式，如果没有users这个表，会自动创建数据表和自动创建表内字段
	//db.AutoMigrate(&User{})

	//使用
	db.Session(&gorm.Session{CreateBatchSize: 100})

	var s string
	s = "lz13526112262@163.com"
	user := User{
		Username: "sdf",
		Password: "15489",
		Email:    &s,
	}
	result := db.Create(&user)

	//打印插入记录的条数
	fmt.Println(result.RowsAffected)
	fmt.Println(result.Error)

	//通过切片插入多条数据
	users := []User{
		{Username: "lsl", Password: "123879", Email: &s},
		{Username: "lsw", Password: "123879", Email: &s},
	}
	result2 := db.Create(&users)
	fmt.Println(result2.RowsAffected)

	//创建记录为指定字段赋值
	create := db.Select("Username", "Email").Create(&user)
	fmt.Println(create.RowsAffected)

	//使用map创建，钩子方法不会执行，关联不会被保存且不会会写主键   无法自动添加create_at 和 update_at
	db.Model(&User{}).Create(map[string]interface{}{
		"Username": "ous", "Password": "542545", "Email": &s,
	})
}
