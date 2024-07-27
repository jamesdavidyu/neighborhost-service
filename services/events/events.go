package events

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/services/auth"
	"github.com/jamesdavidyu/neighborhost-service/utils"
)

type Handler struct {
	store         types.EventStore
	neighborStore types.NeighborStore
	zipcodeStore  types.ZipcodeStore
	addressStore  types.AddressStore
}

func NewHandler(store types.EventStore, neighborStore types.NeighborStore, zipcodeStore types.ZipcodeStore, addressStore types.AddressStore) *Handler {
	return &Handler{store: store, neighborStore: neighborStore, zipcodeStore: zipcodeStore, addressStore: addressStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/events", h.handleGetPublicEvents).Methods("GET")
	router.HandleFunc("/auth/events", auth.WithJWTAuth(h.handleGetZipcodeEvents, h.neighborStore)).Methods("GET")
	router.HandleFunc("/auth/events/neighborhood-events", auth.WithJWTAuth(h.handleGetNeighborhoodEvents, h.neighborStore)).Methods("GET")
	router.HandleFunc("/auth/events/create-event", auth.WithJWTAuth(h.handleCreateEvent, h.neighborStore)).Methods("POST")
}

func (h *Handler) handleGetPublicEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.store.GetPublicEvents()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) handleGetZipcodeEvents(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())

	getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	getZipcodeData, err := h.zipcodeStore.GetZipcodeData(getNeighbor.Zipcode)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	location, err := time.LoadLocation(getZipcodeData.Timezone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	events, err := h.store.GetEventsByZipcode(getNeighbor.Zipcode, time.Now().In(location))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) handleGetNeighborhoodEvents(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())

	getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	getZipcodeData, err := h.zipcodeStore.GetZipcodeData(getNeighbor.Zipcode)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	location, err := time.LoadLocation(getZipcodeData.Timezone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	events, err := h.store.GetEventsByNeighborhoodId(getNeighbor.NeighborhoodId, time.Now().In(location))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())
	var event types.EventPayload

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(event); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
		return
	}

	validateZipcode, err := h.zipcodeStore.ValidateZipcode(
		event.City,
		event.State,
		event.Zipcode,
	)
	if validateZipcode == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
		return
	} else {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
			return
		}

		getAddress, err := h.addressStore.GetAddressByNeighborId(neighborId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
			return
		}

		getAddressId, err := h.addressStore.GetAddressIdByAddress(
			getAddress.FirstName,
			getAddress.LastName,
			event.Address,
			event.City,
			event.State,
			event.Zipcode,
			neighborId,
		)

		if getAddressId.Id == 0 {
			err = h.addressStore.CreateAddress(types.Addresses{
				FirstName:      getAddress.FirstName,
				LastName:       getAddress.LastName,
				Address:        event.Address,
				City:           event.City,
				State:          event.State,
				Zipcode:        event.Zipcode,
				NeighborId:     neighborId,
				NeighborhoodId: getAddress.NeighborhoodId,
			})
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			getNewAddressId, err := h.addressStore.GetAddressIdByAddress(
				getAddress.FirstName,
				getAddress.LastName,
				event.Address,
				event.City,
				event.State,
				event.Zipcode,
				neighborId,
			)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
				return
			}

			err = h.store.CreateEvent(types.Events{
				Name:           event.Name,
				Description:    event.Description,
				Start:          event.Start,
				End:            event.End,
				Reoccurrence:   event.Reoccurrence,
				ForUnloggedins: event.ForUnloggedins,
				ForUnverifieds: event.ForUnverifieds,
				InviteOnly:     event.InviteOnly,
				HostId:         neighborId,
				AddressId:      getNewAddressId.Id,
			})
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(event)
		} else {
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
				return
			}

			err = h.store.CreateEvent(types.Events{
				Name:           event.Name,
				Description:    event.Description,
				Start:          event.Start,
				End:            event.End,
				Reoccurrence:   event.Reoccurrence,
				ForUnloggedins: event.ForUnloggedins,
				ForUnverifieds: event.ForUnverifieds,
				InviteOnly:     event.InviteOnly,
				HostId:         neighborId,
				AddressId:      getAddressId.Id,
			})
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(event)
		}
	}
}
