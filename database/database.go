package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/chrislonng/starx/log"
)

var db *sql.DB


func OpenConnection() {
	var err error
	db, err = sql.Open(
		"mysql",
		"username:password@tcp(192.168.10.10:3306)/bbo")
	if err != nil {
		log.Fatal(err)
	}
}

func CloseConnection() {
	db.Close()
}

func DB() *sql.DB {
	return db
}