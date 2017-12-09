package db

import (
	"database/sql"
	"fmt"
)

func CreateConnection() *sql.DB {
	db, err := sql.Open("postgres", "user=neotek password=kringstone dbname=projectdb sslmode=disable")
    if err != nil {
        fmt.Println(err.Error())
    }
    if err = db.Ping(); err != nil {
        panic(err)
    } else {
        fmt.Println("DB connected")
    }
	return db
}