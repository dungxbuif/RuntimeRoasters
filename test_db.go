package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://user:password@127.0.0.1:54321/auth_db?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
	}
}
