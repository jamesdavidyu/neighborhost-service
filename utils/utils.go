/*
1. GENERAL
2. FOR ZIPCODE CONTROLLERS
3. FOR NEIGHBORS CONTROLLERS
4. FOR ADDRESSES CONTROLLERS
5. FOR NEIGHBORHOODS CONTROLLERS
6. FOR EVENT CONTROLLERS
*/

package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
)

/* 1. GENERAL */

var Validate = validator.New()

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

/* 2. FOR ZIPCODES CONTROLLERS */

func ScanRowIntoZipcodes(rows *sql.Rows) (*types.Zipcodes, error) {
	zipcodeData := new(types.Zipcodes)

	err := rows.Scan(
		&zipcodeData.Zipcode,
		&zipcodeData.City,
		&zipcodeData.State,
		&zipcodeData.Timezone,
	)
	if err != nil {
		return nil, err
	}

	return zipcodeData, nil
}

/* 3. FOR NEIGHBORS CONTROLLERS */

func ScanRowIntoNeighbor(rows *sql.Rows) (*types.Neighbors, error) {
	neighbor := new(types.Neighbors)

	err := rows.Scan(
		&neighbor.Id,
		&neighbor.Email,
		&neighbor.Username,
		&neighbor.Zipcode,
		&neighbor.Password,
		&neighbor.Verified,
		&neighbor.NeighborhoodId,
		&neighbor.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return neighbor, nil
}

/* 4. FOR ADDRESSES CONTROLLERS */

func ScanRowIntoAddresses(rows *sql.Rows) (*types.Addresses, error) {
	addresses := new(types.Addresses)

	err := rows.Scan(
		&addresses.Id,
		&addresses.FirstName,
		&addresses.LastName,
		&addresses.Address,
		&addresses.City,
		&addresses.State,
		&addresses.Zipcode,
		&addresses.NeighborId,
		&addresses.NeighborhoodId,
		&addresses.RecordedAt,
	)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

/* 5. FOR NEIGHBORHOODS CONTROLLERS */

func ScanRowsIntoNeighborhood(rows *sql.Rows) (*types.Neighborhoods, error) {
	neighborhood := new(types.Neighborhoods)

	err := rows.Scan(
		&neighborhood.Id,
		&neighborhood.Neighborhood,
		&neighborhood.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return neighborhood, nil
}

/* 6. FOR EVENT CONTROLLERS */

func ScanRowIntoPublicEvents(rows *sql.Rows) (*types.Events, error) {
	events := new(types.Events)

	err := rows.Scan(
		&events.Id,
		&events.Name,
		&events.Description,
		&events.Start,
		&events.End,
		&events.Reoccurrence,
		&events.ForUnloggedins,
		&events.ForUnverifieds,
		&events.InviteOnly,
		&events.HostId,
		&events.AddressId,
		&events.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func ScanRowIntoNeighborEvents(rows *sql.Rows) (*types.EventAddresses, error) {
	events := new(types.EventAddresses)

	err := rows.Scan(
		&events.Id,
		&events.Name,
		&events.Description,
		&events.Start,
		&events.End,
		&events.Reoccurrence,
		&events.ForUnloggedins,
		&events.ForUnverifieds,
		&events.InviteOnly,
		&events.HostId,
		&events.AddressId,
		&events.CreatedAt,
		&events.AddressAddressId,
		&events.FirstName,
		&events.LastName,
		&events.Address,
		&events.City,
		&events.State,
		&events.Zipcode,
		&events.NeighborId,
		&events.NeighborhoodId,
		&events.RecordedAt,
	)
	if err != nil {
		return nil, err
	}

	return events, nil
}
