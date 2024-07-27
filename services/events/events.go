package events

import (
	"encoding/json"
	"fmt"
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
	router.HandleFunc("/auth/events", auth.WithJWTAuth(h.handleZipcodeEventsWithDate, h.neighborStore)).Methods("POST")
	router.HandleFunc("/auth/events/neighborhood-events", auth.WithJWTAuth(h.handleGetNeighborhoodEvents, h.neighborStore)).Methods("GET")
	router.HandleFunc("/auth/events/create-event", auth.WithJWTAuth(h.handleCreateEvent, h.neighborStore)).Methods("POST")
}

func (h *Handler) handleGetPublicEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.store.GetPublicEvents()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) handleGetZipcodeEvents(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())

	getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	getZipcodeData, err := h.zipcodeStore.GetZipcodeData(getNeighbor.Zipcode)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	location, err := time.LoadLocation(getZipcodeData.Timezone)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	events, err := h.store.GetEventsByZipcode(getNeighbor.Zipcode, time.Now().In(location))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) handleZipcodeEventsWithDate(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())
	var eventFilter types.FilterEventPayload

	if err := json.NewDecoder(r.Body).Decode(&eventFilter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	if err := utils.Validate.Struct(eventFilter); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
		return
	}

	getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	if eventFilter.DateFilter == "On" {
		events, err := h.store.ZipcodeEventsOnDate(getNeighbor.Zipcode, eventFilter.DateTime)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(events)

	} else if eventFilter.DateFilter == "Before" {
		events, err := h.store.ZipcodeEventsBeforeDate(getNeighbor.Zipcode, eventFilter.DateTime)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(events)

	} else if eventFilter.DateFilter == "After" {
		events, err := h.store.ZipcodeEventsAfterDate(getNeighbor.Zipcode, eventFilter.DateTime)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(events)
	}
}

func (h *Handler) handleGetNeighborhoodEvents(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())

	getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	getZipcodeData, err := h.zipcodeStore.GetZipcodeData(getNeighbor.Zipcode)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	location, err := time.LoadLocation(getZipcodeData.Timezone)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	events, err := h.store.GetEventsByNeighborhoodId(getNeighbor.NeighborhoodId, time.Now().In(location))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())
	var event types.CreateEventPayload

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
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

	if validateZipcode.Zipcode == "" {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	} else {
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		getAddress, err := h.addressStore.GetAddressByNeighborId(neighborId)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
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
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("first time using address"))
			return
		} else {
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
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
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(event)
		}
	}
}
