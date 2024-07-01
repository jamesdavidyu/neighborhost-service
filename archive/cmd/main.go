package main

import (
	"database/sql"
	"log"

	// "github.com/jamesdavidyu/neighborhost-service/cmd/api"
	"github.com/jamesdavidyu/neighborhost-service/db"
)

func main() {
	db, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	// server := api.NewAPIServer(":8080", db)
	// if err := server.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
