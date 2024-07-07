package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	neighborhoodControllers "github.com/jamesdavidyu/neighborhost-service/controllers/neighborhoods"
	neighborControllers "github.com/jamesdavidyu/neighborhost-service/controllers/neighbors"
	neighborhoodServices "github.com/jamesdavidyu/neighborhost-service/services/neighborhoods"
	neighborServices "github.com/jamesdavidyu/neighborhost-service/services/neighbors"
	"github.com/jamesdavidyu/neighborhost-service/utils"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

var Port = os.Getenv("PORT")

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/status", getStatus()).Methods("GET")
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	neighborStore := neighborControllers.NewStore(s.db)
	neighborHandler := neighborServices.NewHandler(neighborStore)
	neighborHandler.RegisterRoutes(subrouter)

	neighborhoodStore := neighborhoodControllers.NewStore(s.db)
	neighborhoodHandler := neighborhoodServices.NewHandler(neighborhoodStore, neighborStore)
	neighborhoodHandler.RegisterRoutes(subrouter)

	if Port == "" {
		Port = "8080"

	}

	enhancedRouter := utils.EnableCORS(utils.JSONContentTypeMiddleware(router))

	log.Println("listening on", Port)
	return http.ListenAndServe(":"+Port, enhancedRouter)
}

func getStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		okStatus := map[string]string{"status": "ok"}
		if err := json.NewEncoder(w).Encode(okStatus); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
