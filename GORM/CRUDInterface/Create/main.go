package main

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Email    *string
}

type CreditCard struct {
	gorm.Model
	Number    string
	StudentID uint
}
type Student struct {
	gorm.Model
	Name       string
	CreditCard CreditCard
}

var email = "www.guitu.life"

func main() {
	//连接数据库
	dsn := "root:251210@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	sqlDB, _ := sql.Open("mysql", dsn) //创建数据库连接
	//通过已有的数据库连接来初始化gorm.DB
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}))
	if err != nil {
		log.Fatalln("连接数据库错误")
	}
	//db.AutoMigrate(&Student{})
	//db.AutoMigrate(&CreditCard{})

	//用指定字段添加记录  create_at 和update_at会自动创建
	//var user1 = User{
	//	Username: "qwe",
	//	Password: "123456",
	//}
	//db.Select("username", "password").Create(&user1)

	//批量插入   无法向create传递结构体，只能传递指针类型  要么切片类型是指针类型，要么传递切片地址
	//var users1 = []User{{Username: "wer", Password: "58866", Email: &email}, {Username: "ert", Password: "sss855", Email: &email}}
	//db.Create(&users1)

	//跳过钩子方法创建
	//var user2 = User{
	//	Username: "tyu",
	//	Password: "258512",
	//	Email:    &email,
	//}
	//db.Session(&gorm.Session{SkipHooks: true}).Create(&user2)

	//根据map创建  使用map创建不会自动生成create_at和update_at
	//db.Model(&User{}).Create(map[string]interface{}{
	//	"username": "uio",
	//	"password": "458921",
	//	"email":    &email,
	//})

	//map切片创建
	//db.Model(&User{}).Create([]map[string]interface{}{
	//	{"username": "iop", "password": "dfg", "email": &email},
	//	{"username": "dfg", "password": "ghj", "email": &email},
	//})

	//关联创建
	//db.Create(&Student{Name: "bnm", CreditCard: CreditCard{Number: "41111111111"}})

	//设置默认值
	type Book struct {
		ID   int64
		Name string `gorm:"default:'你当像鸟飞往你的山'"`
	}
	db.AutoMigrate(&Book{})
	db.Create(&Book{Name: "sss"})

	//当字段值是零值时

}
