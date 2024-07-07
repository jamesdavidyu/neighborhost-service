package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/db"
	"github.com/jamesdavidyu/neighborhost-service/routes"
)

// go:embed templates/*
// var resources embed.FS

// var t = template.Must(template.ParseFS(resources, "templates/*"))

var Port = os.Getenv("PORT")

func main() {
	db, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	if Port == "" {
		Port = "8080"

	}

	server := routes.NewAPIServer(":"+Port, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	data := map[string]string{
	// 		"Region": os.Getenv("FLY_REGION"),
	// 	}

	// 	t.ExecuteTemplate(w, "index.html.tmpl", data)
	// })
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
