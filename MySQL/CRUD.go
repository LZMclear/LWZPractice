package main

import "fmt"

type test struct {
	Id       int    `db:"id"`
	Username string `db:"Username"`
	Password string `db:"Password"`
}

func queryRow() {
	sqlStr := "select id, username, password from test where id =?"
	var u test
	err := db.QueryRow(sqlStr, 1).Scan(&u.Id, &u.Username, &u.Password)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id: %d, username: %s, password: %s", u.Id, u.Username, u.Password)
}

func queryMultiRow() {
	sqlStr := "select * from test where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query rows failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	//循环接收查询的数据
	//用一个切片接收
	var users = make([]test, 10)
	i := 0
	for rows.Next() {
		err := rows.Scan(&users[i].Id, &users[i].Username, &users[i].Password)
		if err != nil {
			fmt.Printf("scan failed! err:%v\n", err)
		}
		i++
	}
	fmt.Println(users)
}

func insertRow() {
	sqlStr := "insert into test(username,password) values (?, ?)"
	result, err := db.Exec(sqlStr, "归途", "123456")
	if err != nil {
		fmt.Printf("insert failed! err:%v\n", err)
		return
	}
	theID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert id failed! err:%v\n", err)
	}
	fmt.Printf("insert success, the id is %d\n", theID)
}

func updateRow() {
	sqlStr := "update test set password=? where id=?"
	result, err := db.Exec(sqlStr, "456123", 1)
	if err != nil {
		fmt.Printf("update row failed! err:%v\n", err)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("get affectrow failed, err:%v\n ", err)
		return
	}
	fmt.Printf("update row success! affected row is: %d", n)
}

func deleteRow() {
	sqlStr := "delete from test where id=?"
	result, err := db.Exec(sqlStr, 3)
	if err != nil {
		fmt.Printf("delete row failed! err:%v\n", err)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("get affected row failed! err:%v\n", err)
		return
	}
	fmt.Printf("第%d行被成功删除", n)
}

func transactionDemo() {
	tx, err := db.Begin() //开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "Update user set age=40 where id=?"
	ret2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql2 failed, err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}
	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}

	fmt.Println("exec trans success!")
}
