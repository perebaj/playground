package main

import (
	"database/sql"
)

func OpenConn() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = db.Ping()

	return db, err
}

func main() {
	db, err := OpenConn()
	if err != nil {
		panic(err)
	}
	defer db.Close()
}