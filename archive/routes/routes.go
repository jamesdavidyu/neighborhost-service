package routes

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"
// 	_ "github.com/jackc/pgx/v4/stdlib"
// 	"github.com/jamesdavidyu/neighborhost-service/cmd/model/db"
// 	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
// 	"github.com/jamesdavidyu/neighborhost-service/controllers"
// 	"github.com/jamesdavidyu/neighborhost-service/services"
// )

// type Handler struct {
// 	store types.NeighborStore
// }

// func NewHandler(store types.NeighborStore) *Handler {
// 	return &Handler{store: store}
// }

// func (h *Handler) Routes(router *mux.Router) {
// 	db, err := db.DB()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	router.HandleFunc("/status", GetStatus()).Methods("GET")
// 	router.HandleFunc("/neighborhoods", controllers.CreateNeighborhood(db)).Methods("POST")
// 	router.HandleFunc("/auth/register", services.CreateNeighbor(db)).Methods("POST")
// 	router.HandleFunc("/auth/login", services.Login(db)).Methods("POST")
// }

// func GetStatus() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		okStatus := map[string]string{"status": "ok"}
// 		if err := json.NewEncoder(w).Encode(okStatus); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// }
