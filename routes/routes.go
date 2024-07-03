package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/db"
	"github.com/jamesdavidyu/neighborhost-service/controllers"
	"github.com/jamesdavidyu/neighborhost-service/utils"
)

func Routes() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	db, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	routePrefix := "/api/v1"
	router.HandleFunc(routePrefix+"/status", GetStatus()).Methods("GET")
	router.HandleFunc(routePrefix+"/neighborhoods", controllers.GetNeighborhoods(db)).Methods("GET")
	router.HandleFunc(routePrefix+"/neighborhoods", controllers.CreateNeighborhood(db)).Methods("POST")
	router.HandleFunc(routePrefix+"/auth/register", controllers.CreateNeighbor(db)).Methods("POST")

	enhancedRouter := utils.EnableCORS(utils.JSONContentTypeMiddleware(router))

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, enhancedRouter))
}

func GetStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		okStatus := map[string]string{"status": "ok"}
		if err := json.NewEncoder(w).Encode(okStatus); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
