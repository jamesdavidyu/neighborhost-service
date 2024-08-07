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
	router.HandleFunc("/events/auth", auth.WithJWTAuth(h.handleGetEvents, h.neighborStore)).Methods("GET")
	// router.HandleFunc("/events/location-filter/auth", auth.WithJWTAuth(h.handleEventsWithLocation, h.neighborStore)).Methods("POST") // need to figure out how null filters look before writing this endpoint?
	// router.HandleFunc("/events/date-filter/auth", auth.WithJWTAuth(h.handleZipcodeEventsWithDate, h.neighborStore)).Methods("POST")
	// router.HandleFunc("/events/location-date-filter/auth", auth.WithJWTAuth(h.handleLocationEventsWithDate, h.neighborStore)).Methods("POST")
	router.HandleFunc("/events/create-event/auth", auth.WithJWTAuth(h.handleCreateEvent, h.neighborStore)).Methods("POST")
}

func (h *Handler) handleGetPublicEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.store.GetPublicEvents()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	neighborId := auth.GetNeighborIdFromContext(r.Context())
	var eventFilters types.EventFilterPayload

	getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	getZipcodeData, err := h.zipcodeStore.GetZipcodeData(getNeighbor.Zipcode)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	location, err := time.LoadLocation(getZipcodeData.Timezone)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	qs := r.URL.Query()
	eventFilters.LocationFilter = utils.ReadString(qs, "location-filter", "My zipcode")
	eventFilters.DateFilter = utils.ReadString(qs, "date-filter", "") // default is on after depending on how client makes request
	eventFilters.DateTime = utils.ReadDateTime(qs, "datetime", time.Now().In(location))
	// build out controllers

	if eventFilters.LocationFilter == "My zipcode" {
		if eventFilters.DateFilter == "On" {
			events, err := h.store.GetZipcodeEventsOnDate(getNeighbor.Zipcode, eventFilters.DateTime.In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		} else if eventFilters.DateFilter == "Before" {
			events, err := h.store.GetZipcodeEventsBeforeDate(getNeighbor.Zipcode, eventFilters.DateTime.In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		} else if eventFilters.DateFilter == "After" {
			events, err := h.store.GetZipcodeEventsAfterDate(getNeighbor.Zipcode, eventFilters.DateTime.In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

			// TODO: two more controllers for on/after and on/before

		} else {
			events, err := h.store.GetEventsByZipcode(getNeighbor.Zipcode, time.Now().In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		}

	} else if eventFilters.LocationFilter == "My neighborhood" {
		if eventFilters.DateFilter == "On" {
			events, err := h.store.GetNeighborhoodEventsOnDate(getNeighbor.NeighborhoodId, eventFilters.DateTime.In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		} else if eventFilters.DateFilter == "Before" {
			events, err := h.store.GetNeighborhoodEventsBeforeDate(getNeighbor.NeighborhoodId, eventFilters.DateTime.In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		} else if eventFilters.DateFilter == "After" {
			events, err := h.store.GetNeighborhoodEventsAfterDate(getNeighbor.NeighborhoodId, eventFilters.DateTime.In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		} else {
			events, err := h.store.GetEventsByNeighborhoodId(getNeighbor.NeighborhoodId, time.Now().In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		}

	} else if eventFilters.LocationFilter == "My city" {
		getAddress, err := h.addressStore.GetAddressByNeighborId(neighborId) // this will need to get fixed to specify home address
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		if eventFilters.DateFilter == "On" {
			events, err := h.store.GetCityEventsOnDate(getAddress.City, getAddress.State, getAddress.Zipcode, eventFilters.DateTime.In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		} else if eventFilters.DateFilter == "Before" {
			events, err := h.store.GetCityEventsBeforeDate(getAddress.City, getAddress.State, getAddress.Zipcode, eventFilters.DateTime.In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		} else if eventFilters.DateFilter == "After" {
			events, err := h.store.GetCityEventsAfterDate(getAddress.City, getAddress.State, getAddress.Zipcode, eventFilters.DateTime.In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		} else {
			events, err := h.store.GetEventsByCity(getAddress.City, getAddress.State, getAddress.Zipcode, time.Now().In(location))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
				return
			}

			utils.WriteJSON(w, http.StatusOK, events)

		}

	} else {
		// TODO: create controller and logic for if user inputs a city
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}
}

// func (h *Handler) handleZipcodeEventsWithDate(w http.ResponseWriter, r *http.Request) {
// 	neighborId := auth.GetNeighborIdFromContext(r.Context())
// 	var dateFilter types.DateFilterPayload

// 	if err := json.NewDecoder(r.Body).Decode(&dateFilter); err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
// 		return
// 	}

// 	if err := utils.Validate.Struct(dateFilter); err != nil {
// 		errors := err.(validator.ValidationErrors)
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
// 		return
// 	}

// 	getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
// 		return
// 	}

// 	getZipcodeData, err := h.zipcodeStore.GetZipcodeData(getNeighbor.Zipcode)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
// 		return
// 	}

// 	location, err := time.LoadLocation(getZipcodeData.Timezone)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
// 		return
// 	}

// 	if dateFilter.DateFilter == "On" {

// 	} else if dateFilter.DateFilter == "Before" {
// 		events, err := h.store.GetZipcodeEventsBeforeDate(getNeighbor.Zipcode, dateFilter.DateTime.In(location))
// 		if err != nil {
// 			utils.WriteError(w, http.StatusInternalServerError, err)
// 			return
// 		}

// 		utils.WriteJSON(w, http.StatusOK, events)

// 	} else if dateFilter.DateFilter == "After" {
// 		events, err := h.store.GetZipcodeEventsAfterDate(getNeighbor.Zipcode, dateFilter.DateTime.In(location))
// 		if err != nil {
// 			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
// 			return
// 		}

// 		utils.WriteJSON(w, http.StatusOK, events)

// 	} else {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
// 		return
// 	}
// }

// func (h *Handler) handleEventsWithLocation(w http.ResponseWriter, r *http.Request) {
// 	neighborId := auth.GetNeighborIdFromContext(r.Context())
// 	var locationFilter types.LocationFilterPayload

// 	if err := json.NewDecoder(r.Body).Decode(&locationFilter); err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
// 		return
// 	}

// 	if err := utils.Validate.Struct(locationFilter); err != nil {
// 		errors := err.(validator.ValidationErrors)
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
// 		return
// 	}

// 	getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
// 		return
// 	}

// 	getZipcodeData, err := h.zipcodeStore.GetZipcodeData(getNeighbor.Zipcode)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
// 		return
// 	}

// 	location, err := time.LoadLocation(getZipcodeData.Timezone)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
// 		return
// 	}

// 	if locationFilter.LocationFilter == "My zipcode" {
// 		h.handleGetZipcodeEvents(w, r)

// 	} else if locationFilter.LocationFilter == "My neighborhood" {

// 	} else if locationFilter.LocationFilter == "All" {
// 		// events, err := h.store.GetAllEvents(time.Now().In(location))
// 		// if err != nil {
// 		// 	utils.WriteError(w, http.StatusInternalServerError, err)
// 		// 	return
// 		// }

// 		// utils.WriteJSON(w, http.StatusOK, events)

// 	}
// }

// func (h *Handler) handleLocationEventsWithDate(w http.ResponseWriter, r *http.Request) {
// 	neighborId := auth.GetNeighborIdFromContext(r.Context())
// 	var locationDateFilter types.LocationDateFilterPayload

// 	if err := json.NewDecoder(r.Body).Decode(&locationDateFilter); err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
// 		return
// 	}

// 	if err := utils.Validate.Struct(locationDateFilter); err != nil {
// 		errors := err.(validator.ValidationErrors)
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
// 		return
// 	}

// 	getNeighbor, err := h.neighborStore.GetNeighborById(neighborId)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
// 		return
// 	}

// 	getZipcodeData, err := h.zipcodeStore.GetZipcodeData(getNeighbor.Zipcode)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
// 		return
// 	}

// 	location, err := time.LoadLocation(getZipcodeData.Timezone)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found"))
// 		return
// 	}

// 	if locationDateFilter.LocationFilter == "My zipcode" {
// 		h.handleZipcodeEventsWithDate(w, r)

// 	} else if locationDateFilter.LocationFilter == "My neighborhood" {
// 		 else {
// 			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
// 			return
// 		}

// 	} else if locationDateFilter.LocationFilter == "My city" {
// 		getAddress, err := h.addressStore.GetAddressByNeighborId(neighborId)
// 		if err != nil {
// 			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
// 			return
// 		}

// 		 else {
// 			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
// 			return
// 		}
// 	} else {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
// 		return
// 	}
// }

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

	checkAddress, err := h.addressStore.GetAddressIdByAddress(
		event.Address,
		event.City,
		event.State,
		event.Zipcode,
		event.Type,
		neighborId,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// TODO: need to add logic where if address id is zero doesn't exist, it's created

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

	if checkAddress.Id == 0 {
		getHomeAddress, err := h.addressStore.GetAddressByNeighborId(getNeighbor.Id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		err = h.addressStore.CreateAddress(types.Addresses{
			FirstName:      getHomeAddress.FirstName,
			LastName:       getHomeAddress.LastName,
			Address:        event.Address,
			City:           event.City,
			State:          event.State,
			Zipcode:        event.Zipcode,
			Type:           event.Type,
			NeighborId:     neighborId,
			NeighborhoodId: 1, // need to delete this line later
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		checkAddressAgain, err := h.addressStore.GetAddressIdByAddress(
			event.Address,
			event.City,
			event.State,
			event.Zipcode,
			event.Type,
			neighborId,
		)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		err = h.store.CreateEvent(types.Events{
			Name:           event.Name,
			Description:    event.Description,
			Start:          event.Start.In(location),
			End:            event.End.In(location),
			Reoccurrence:   event.Reoccurrence,
			ForUnloggedins: event.ForUnloggedins,
			ForUnverifieds: event.ForUnverifieds,
			InviteOnly:     event.InviteOnly,
			HostId:         neighborId,
			AddressId:      checkAddressAgain.Id,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(event)

	} else {
		err = h.store.CreateEvent(types.Events{
			Name:           event.Name,
			Description:    event.Description,
			Start:          event.Start.In(location),
			End:            event.End.In(location),
			Reoccurrence:   event.Reoccurrence,
			ForUnloggedins: event.ForUnloggedins,
			ForUnverifieds: event.ForUnverifieds,
			InviteOnly:     event.InviteOnly,
			HostId:         neighborId,
			AddressId:      checkAddress.Id,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(event)
	}

}
