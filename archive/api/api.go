package api

// import (
// 	"database/sql"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"
// 	"github.com/jamesdavidyu/neighborhost-service/utils"
// )

// type APIServer struct {
// 	addr string
// 	db   *sql.DB
// }

// func NewAPIServer(addr string, db *sql.DB) *APIServer {
// 	return &APIServer{
// 		addr: addr,
// 		db:   db,
// 	}
// }

// func (s *APIServer) Run() error {
// 	router := mux.NewRouter()
// 	// subrouter := router.PathPrefix("api/v1").Subrouter()

// 	// router.HandleFunc("api/v1/test", service.GetTests(db)).Methods("GET")

// 	log.Println("Listening on", s.addr)

// 	enhancedRouter := utils.EnableCORS(utils.JSONContentTypeMiddleware(router))

// 	return http.ListenAndServe(s.addr, enhancedRouter)
// }

// func enableCORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func jsonContentTypeMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		next.ServeHTTP(w, r)
// 	})
// }
