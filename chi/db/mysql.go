package db

import (
	"fmt"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	connString string
	conn *sql.DB
}

func NewMysql(connString string) *sql.DB {
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	return conn
}