/*
1. GENERAL
2. FOR ZIPCODE CONTROLLERS
3. FOR NEIGHBORS CONTROLLERS/SERVICES
4. FOR ADDRESSES CONTROLLERS
5. FOR NEIGHBORHOODS CONTROLLERS
6. FOR EVENT CONTROLLERS
7. FOR FRIENDS CONTROLLERS
*/

package utils

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

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

func ReadString(qs url.Values, key string, defaultValue string) string {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}

	return strings.Split(value, ",")
}

func ReadInt(qs url.Values, key string, defaultValue int) int {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func ReadDateTime(qs url.Values, key string, defaultValue time.Time) time.Time {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}

	layout := "2006-01-02 15:04:05"

	dateTimeValue, err := time.Parse(layout, value)
	if err != nil {
		return defaultValue
	}

	return dateTimeValue
}

func ToProperCase(input string) string {
	words := strings.Fields(input)
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		words[i] = string(runes)
	}

	return strings.Join(words, " ")
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

/* 3. FOR NEIGHBORS CONTROLLERS/SERVICES */

func ScanRowIntoNeighbor(rows *sql.Rows) (*types.Neighbors, error) {
	neighbor := new(types.Neighbors)

	err := rows.Scan(
		&neighbor.Id,
		&neighbor.Email,
		&neighbor.Username,
		&neighbor.Zipcode,
		&neighbor.Password,
		&neighbor.Verified,
		&neighbor.Ip,
		&neighbor.NeighborhoodId,
		&neighbor.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return neighbor, nil
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
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
		&addresses.Type,
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
		&events.Type,
		&events.NeighborId,
		&events.NeighborhoodId,
		&events.RecordedAt,
	)
	if err != nil {
		return nil, err
	}

	return events, nil
}

/* 7. FOR FRIENDS CONTROLLERS */

func ScanRowIntoFriendsList(rows *sql.Rows) (*types.FriendsList, error) {
	friends := new(types.FriendsList)

	err := rows.Scan(
		&friends.Id,
		&friends.NeighborId,
		&friends.NeighborsFriendId,
		&friends.FriendedAt,
		&friends.NeighborsId,
		&friends.Email,
		&friends.Username,
		&friends.NeighborZipcode,
		&friends.Password,
		&friends.Verified,
		&friends.Ip,
		&friends.NeighborNeighborhoodId,
		&friends.CreatedAt,
		&friends.AddressesId,
		&friends.FirstName,
		&friends.LastName,
		&friends.Address,
		&friends.City,
		&friends.State,
		&friends.Zipcode,
		&friends.Type,
		&friends.AddressesNeighborId,
		&friends.NeighborhoodId,
		&friends.RecordedAt,
	)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func ScanRowIntoFriendRequests(rows *sql.Rows) (*types.PendingFriendRequests, error) {
	friends := new(types.PendingFriendRequests)

	err := rows.Scan(
		&friends.FriendRequestId,
		&friends.NeighborId,
		&friends.RequestedFriendId,
		&friends.Status,
		&friends.FriendRequestedAt,
		&friends.Id,
		&friends.Email,
		&friends.Username,
		&friends.Zipcode,
		&friends.Password,
		&friends.Verified,
		&friends.Ip,
		&friends.NeighborhoodId,
		&friends.CreatedAt,
	) // need all these rows?
	if err != nil {
		return nil, err
	}

	return friends, nil
}

/* FOR PROFILES CONTROLLERS */

func ScanRowIntoProfiles(rows *sql.Rows) (*types.Profiles, error) {
	profiles := new(types.Profiles)

	err := rows.Scan(
		&profiles.NeighborId,
		&profiles.Bio,
		&profiles.DateOfBirth,
		&profiles.DateOfBirthPublic,
		&profiles.Gender,
		&profiles.GenderPublic,
		&profiles.Race,
		&profiles.RacePublic,
		&profiles.Ethnicity,
		&profiles.EthnicityPublic,
		&profiles.RelationshipStatus,
		&profiles.RelationshipStatusPublic,
		&profiles.Religion,
		&profiles.ReligionPublic,
		&profiles.Politics,
		&profiles.PoliticsPublic,
	) // if querying for education, occupation, and interest separately, keep profile information separate?
	if err != nil {
		return nil, err
	}

	return profiles, nil
}
