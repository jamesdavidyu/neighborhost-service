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
	zipcodeStore  types.ZipcodeStore
}

func NewHandler(store types.AddressStore, neighborStore types.NeighborStore, zipcodeStore types.ZipcodeStore) *Handler {
	return &Handler{store: store, neighborStore: neighborStore, zipcodeStore: zipcodeStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/addresses/auth", auth.WithJWTAuth(h.handleGetAddresses, h.neighborStore)).Methods("GET")
	router.HandleFunc("/address/auth", auth.WithJWTAuth(h.handleCreateAddress, h.neighborStore)).Methods("POST")
}

func (h *Handler) handleGetAddresses(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())
	addresses, err := h.store.GetAddressesByNeighborId(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, addresses)
}

func (h *Handler) handleCreateAddress(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())
	var address types.AddressPayload

	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	if err := utils.Validate.Struct(address); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
		return
	}

	validateZipcode, err := h.zipcodeStore.ValidateZipcode(
		address.City,
		address.State,
		address.Zipcode,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	if validateZipcode.Zipcode == "" {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("please double check your city/state/zipcode combination"))
		return
	} else {
		getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		if address.Zipcode != getNeighbor.Zipcode {
			err = h.neighborStore.UpdateZipcodeWithId(types.Neighbors{
				Zipcode: address.Zipcode,
				Id:      neighborId,
			})
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
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
			Type:           address.Type,
			NeighborId:     neighborId,
			NeighborhoodId: 1,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(address) // does this need to return anything?
	}
}
