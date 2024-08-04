package friends

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/services/auth"
	"github.com/jamesdavidyu/neighborhost-service/utils"
)

type Handler struct {
	store         types.FriendStore
	neighborStore types.NeighborStore
}

func NewHandler(store types.FriendStore, neighborStore types.NeighborStore) *Handler {
	return &Handler{store: store, neighborStore: neighborStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/friends/auth", auth.WithJWTAuth(h.handleGetFriends, h.neighborStore)).Methods("GET")
	router.HandleFunc("/friend-requests/{requestingFriendId}/auth", auth.WithJWTAuth(h.handleCreateFriendRequest, h.neighborStore)).Methods("POST")
	router.HandleFunc("/friend-requests/auth", auth.WithJWTAuth(h.handleGetFriendRequests, h.neighborStore)).Methods("GET")
}

func (h *Handler) handleGetFriends(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())

	friends, err := h.store.GetFriendsByNeighborId(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, friends)
}

func (h *Handler) handleCreateFriendRequest(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())
	vars := mux.Vars(r)
	str, ok := vars["requestingFriendId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	requestingFriendId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	err = h.store.CreateFriendRequest(types.FriendRequests{
		NeighborId:         neighborId,
		RequestingFriendId: requestingFriendId,
		Status:             "pending",
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(requestingFriendId)
}

func (h *Handler) handleGetFriendRequests(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())

	friendRequests, err := h.store.GetFriendRequestsByNeighborId(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, friendRequests)
}
