package main

import (
	"database/sql"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, err := sql.Open("postgres", "postgres://test:1234@localhost/test?sslmode=disable")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}

}
