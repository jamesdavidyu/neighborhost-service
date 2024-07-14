package addresses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/services/auth"
	"github.com/jamesdavidyu/neighborhost-service/utils"
)

type Handler struct {
	store         types.AddressStore
	neighborStore types.NeighborStore
}

func NewHandler(store types.AddressStore, neighborStore types.NeighborStore) *Handler {
	return &Handler{store: store, neighborStore: neighborStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/address", auth.WithJWTAuth(h.handleCreateAddress, h.neighborStore)).Methods("POST")
}

func (h *Handler) handleCreateAddress(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())

	var address types.AddressPayload
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(address); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
		return
	}

	getZipcode, err := h.neighborStore.GetNeighborById(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
		return
	}

	if address.Zipcode != getZipcode.Zipcode {
		err = h.neighborStore.UpdateZipcodeWithId(types.Neighbors{
			Zipcode: address.Zipcode,
			Id:      neighborId,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}

	err = h.store.CreateAddress(types.Addresses{
		FirstName:      address.FirstName,
		LastName:       address.LastName,
		Address:        address.Address,
		City:           address.City,
		State:          address.State,
		Zipcode:        address.Zipcode,
		NeighborId:     neighborId,
		NeighborhoodId: getZipcode.NeighborhoodId,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address) // does this need to return anything?
}
