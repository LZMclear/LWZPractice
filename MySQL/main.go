package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 定义一个全局对象db
var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:251210@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	//err := initDB() // 调用输出化数据库的函数
	//if err != nil {
	//	fmt.Printf("init db failed,err:%v\n", err)
	//	return
	//}
	//单行查询
	//queryRow()
	//多行查询
	//queryMultiRow()
	//插入数据
	//insertRow()
	//更新数据
	//updateRow()
	//删除数据
	//deleteRow()

	//sqlx
	err := initSqlxDB()
	if err != nil {
		fmt.Printf("connect db failed! err: %v\n", err)
	}
	//查询
	//sqlxQueryRows()
	//插入
	//insertSqlxRow()
	//更新
	//updateSqlxRow()
	//删除
	//deleteSqlxRow()
	//使用NamedExec插入数据
	//insertNamedExec()
	//事务操作
	//e := sqlxTransaction()
	//fmt.Printf("err:%v\n", e)
	u1 := test{
		Username: "npy",
		Password: "123456",
	}
	u2 := test{
		Username: "wwy",
		Password: "123456",
	}
	u3 := test{
		Username: "gqw",
		Password: "123456",
	}
	users := []interface{}{u1, u2, u3}
	//users := []*test{&u1, &u2, &u3}
	BatchInsertUser(users)
	//BatchInsertUsers3(users)

}
