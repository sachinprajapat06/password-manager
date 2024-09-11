package db

import (
	"database/sql"
	"log"
)

func GetDBConn() *sql.DB {
	db, err := sql.Open("mysql", "root:12345@tcp(127.0.0.1:3306)/passDB")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
