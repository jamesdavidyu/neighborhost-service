package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jamesdavidyu/neighborhost-service/services"
	"github.com/jamesdavidyu/neighborhost-service/utils"
	"github.com/joho/godotenv"
)

func Routes() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	godotenv.Load()

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/neighborhoods", services.GetNeighborhoods(db)).Methods("GET")

	enhancedRouter := utils.EnableCORS(utils.JSONContentTypeMiddleware(router))

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, enhancedRouter))
}
