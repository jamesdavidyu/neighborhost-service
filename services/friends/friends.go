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
	router.HandleFunc("/friend-requests/{requestedFriendId}/auth", auth.WithJWTAuth(h.handleCreateFriendRequest, h.neighborStore)).Methods("POST")
	router.HandleFunc("/friend-requests/auth", auth.WithJWTAuth(h.handleGetFriendRequests, h.neighborStore)).Methods("GET")
	router.HandleFunc("/friend-requests/{friendId}/{status}/auth", auth.WithJWTAuth(h.handlePutFriendRequest, h.neighborStore)).Methods("PUT")
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
	str, ok := vars["requestedFriendId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	requestedFriendId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	err = h.store.CreateFriendRequest(types.FriendRequests{
		NeighborId:        neighborId,
		RequestedFriendId: requestedFriendId,
		Status:            "pending",
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	okStatus := map[string]string{"requestedFriendId": str}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(okStatus)
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

func (h *Handler) handlePutFriendRequest(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())
	vars := mux.Vars(r)
	str, ok := vars["friendId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}
	str1, ok := vars["status"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	friendId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	if str1 == "accepted" {
		err = h.store.UpdateFriendRequest(types.FriendRequests{
			NeighborId:        friendId,
			RequestedFriendId: neighborId,
			Status:            str1,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		err = h.store.CreateFriend(types.Friends{
			NeighborId:        neighborId,
			NeighborsFriendId: friendId,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		okStatus := map[string]string{"friendId": str}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(okStatus)

	} else if str1 == "declined" {
		err = h.store.UpdateFriendRequest(types.FriendRequests{
			NeighborId:        friendId,
			RequestedFriendId: neighborId,
			Status:            str1,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		okStatus := map[string]string{"friendId": str}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(okStatus)
	} else {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}
}
