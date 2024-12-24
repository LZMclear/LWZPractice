package main

import (
	"fmt"
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

func main() {
	//连接数据库
	db, err := gorm.Open(mysql.Open("root:251210@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatalln("open mysql database failed")
		return
	}
	//查询第一条数据到user1中
	var user1 User
	db.First(&user1)

	//获取一条数据，没有指定排序字段
	var user2 User
	db.Take(&user2)

	//获取最后一条记录
	var user3 User
	db.Last(&user3)

	//将查询结果保存到map中
	user4 := map[string]interface{}{}
	db.Table("users").Model(&User{}).First(&user4) //.Table指明要查询数据库中的具体数据表

	//如果主键是数字类型，可以使用内联条件查询  当使用字符串时，应额外避免SQL注入问题
	var user5 User
	db.First(&user5, 10) //10可以使用"10",此时应注意SQL注入问题

	//当目标对象有一个主键值时，将使用主键构成查询条件   结构体内的值不是主键值会出错
	//如果使用了gorm的特定字段，如DeletedAt,那么查询条件会加上DeletedAt不为空
	var user6 = User{Username: "lsl"}
	db.First(&user6)

	//检索全部对象  创建切片用于接收
	var users1 []User
	db.Find(&users1)

	//string条件查询
	var user7 User
	db.Where("username = ?", "osg").First(&user7)

	//struct map条件
	//当使用结构体查询时，gorm只会查询非零的字段，如果字段设置为零值，它不会构成查询条件
	var user8 User
	db.Where(&User{Username: "sdf"}).First(&user8)

	var user9 User
	db.Where(map[string]interface{}{"username": "sdf"}).First(&user9)

	//主键切片
	var users2 []User
	db.Where([]int{1, 2, 3}).Find(&users2)

	//Not条件
	var users3 []User
	db.Not(map[string]interface{}{"username": []string{"sdf", "lsw"}}).Find(&users3) //Not In

	//Or 条件
	var users4 []User
	db.Where("username = ?", "osg").Or(User{Username: "lsl"}).Find(&users4)

	//使用Select选择指定字段
	var users5 []User
	db.Select("username").Find(&users5)

	//Limited指定要检索最大记录数，offset指定返回记录时跳过的记录数
	var users6 []User
	db.Limit(3).Offset(3).Find(&users6)

	//distinct选取不同的值
	var users7 []User
	db.Distinct("username").Order("username").Find(&users7)
	fmt.Println(users7)
}
