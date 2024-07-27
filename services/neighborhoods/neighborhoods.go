package neighborhoods

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/services/auth"
	"github.com/jamesdavidyu/neighborhost-service/utils"
)

type Handler struct {
	store         types.NeighborhoodStore
	neighborStore types.NeighborStore
}

func NewHandler(store types.NeighborhoodStore, neighborStore types.NeighborStore) *Handler {
	return &Handler{store: store, neighborStore: neighborStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/neighborhoods", auth.WithJWTAuth(h.handleGetNeighborhoods, h.neighborStore)).Methods("GET")
}

func (h *Handler) handleGetNeighborhoods(w http.ResponseWriter, r *http.Request) {
	neighborhoods, err := h.store.GetNeighborhoods()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, neighborhoods)
}
