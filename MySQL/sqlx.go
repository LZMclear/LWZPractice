package main

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var sqlxdb *sqlx.DB

func initSqlxDB() (err error) {
	dsn := "root:251210@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	sqlxdb, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	//设置数据库最大连接数目
	sqlxdb.SetMaxOpenConns(20)
	//设置数据库最大闲置连接数
	sqlxdb.SetMaxIdleConns(10)
	return
}

func sqlxQueryRow() {
	sqlStr := "select * from test where id =?"
	var u test
	err := sqlxdb.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d username:%s password:%s\n", u.Id, u.Username, u.Password)
}

func sqlxQueryRows() {
	sqlStr := "select * from test where id > ?"
	var users []test
	err := sqlxdb.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", users)
}

func insertSqlxRow() {
	sqlxStr := "insert into test(username, password) values (?,?)"
	result, err := sqlxdb.Exec(sqlxStr, "gsy", "1014")
	if err != nil {
		fmt.Printf("insert data failed! err: %v\n", err)
		return
	}
	ID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("get last id failed! err: %v\n", err)
		return
	}
	fmt.Printf("insert data success, the last id is %d\n", ID)
}

func updateSqlxRow() {
	sqlxStr := "update test set password = ? where id =?"
	result, err := sqlxdb.Exec(sqlxStr, "123456", 1)
	if err != nil {
		fmt.Printf("update data failed! err: %v\n", err)
		return
	}
	id, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("get affected row failed! err: %v\n", err)
		return
	}
	fmt.Printf("update data success, update row's id is %d\n ", id)
}

func deleteSqlxRow() {
	sqlxStr := "delete from test where id = ?"
	result, err := sqlxdb.Exec(sqlxStr, 3)
	if err != nil {
		fmt.Printf("delete data failed! err: %v\n", err)
		return
	}
	id, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("get affected row failed! err: %v\n", err)
		return
	}
	fmt.Printf("delete row success,delete id is %d", id)
}

// NamedExec示例
func insertNamedExec() {
	sqlxStr := "insert into test (username) values (:name)"
	_, err := sqlxdb.NamedExec(sqlxStr,
		map[string]interface{}{
			"name": "归途"})
	if err != nil {
		fmt.Printf("insert data failed! err: %v\n", err)
		return
	}
	fmt.Println("operation success!")
}

// 事务操作
func sqlxTransaction() (err error) {
	tx, err := sqlxdb.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()

	sqlStr1 := "Update test set username=? where id= ?"
	rs, err := tx.Exec(sqlStr1, "ldl", 1)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	sqlStr2 := "Update test set password=? where id=?"
	rs, err = tx.Exec(sqlStr2, "789456", 4)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	return err
}
