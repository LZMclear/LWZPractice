package utils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDB() (db *sqlx.DB) {
	dsn := "root:251210@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	return db
}
