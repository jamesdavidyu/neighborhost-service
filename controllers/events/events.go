/*
1. PUBLIC
2. ZIPCODE
3. NEIGHBORHOOD
4. CITY
5. GENERAL
*/

package events

import (
	"database/sql"
	"time"

	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

/* 1. PUBLIC */

func (s *Store) GetPublicEvents() ([]types.Events, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events
		WHERE for_unloggedins = TRUE
		AND start >= $1
		ORDER BY start
		LIMIT 10`, time.Now(), // AB test to figure out optimal limit number
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.Events, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoPublicEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

/* 2. ZIPCODE */

func (s *Store) GetEventsByZipcode(zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1
		AND start >= $2
		ORDER BY start`, zipcode, dateTime,
	) // need to add unloggedin, unverified, inviteOnly filters
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetZipcodeEventsOnDate(zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1
		AND start = $2
		ORDER BY start`, zipcode, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetZipcodeEventsBeforeDate(zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1
		AND start < $2
		ORDER BY start DESC`, zipcode, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetZipcodeEventsAfterDate(zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1
		AND start > $2
		ORDER BY start`, zipcode, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

// need to add controllers for on/before and on/after

/* 3. NEIGHBORHOOD */

func (s *Store) GetEventsByNeighborhoodId(neighborhoodId int, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.neighborhood_id = $1
		AND start >= $2
		ORDER BY start`, neighborhoodId, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetNeighborhoodEventsOnDate(neighborhoodId int, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e 
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.neighborhood_id = $1
		AND start = $2
		ORDER BY start`, neighborhoodId, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetNeighborhoodEventsBeforeDate(neighborhoodId int, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e 
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.neighborhood_id = $1
		AND start < $2
		ORDER BY start DESC`, neighborhoodId, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetNeighborhoodEventsAfterDate(neighborhoodId int, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e 
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.neighborhood_id = $1
		AND start > $2
		ORDER BY start`, neighborhoodId, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

/* 4. CITY */

func (s *Store) GetEventsByCity(city string, state string, zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.city = $1 AND a.state = $2 AND a.zipcode = $3
		AND start >= $4
		ORDER BY start`, city, state, zipcode, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetCityEventsOnDate(city string, state string, zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.city = $1 AND a.state = $2 AND a.zipcode = $3
		AND start = $4
		ORDER BY start`, city, state, zipcode, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetCityEventsBeforeDate(city string, state string, zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.city = $1 AND a.state = $2 AND a.zipcode = $3
		AND start < $4
		ORDER BY start DESC`, city, state, zipcode, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetCityEventsAfterDate(city string, state string, zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.city = $1 AND a.state = $2 AND a.zipcode = $3
		AND start > $4
		ORDER BY start`, city, state, zipcode, dateTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := utils.ScanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

/* 5. ALL, needed? */

// func (s *Store) GetAllEvents(dateTime time.Time) ([]types.EventAddresses, error) {
// 	rows, err := s.db.Query(
// 		`SELECT * FROM events e
// 		LEFT OUTER JOIN addresses a ON a.id = e.address_id
// 		WHERE start >= $1`, dateTime,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	events := make([]types.EventAddresses, 0)
// 	for rows.Next() {
// 		event, err := utils.ScanRowIntoNeighborEvents(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		events = append(events, *event)
// 	}

// 	return events, nil
// }

/* 6. GENERAL */

func (s *Store) CreateEvent(event types.Events) error {
	_, err := s.db.Exec(
		`INSERT INTO events (
			name,
			description,
			start,
			"end",
			reoccurrence,
			for_unloggedins,
			for_unverifieds,
			invite_only,
			host_id,
			address_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		event.Name,
		event.Description,
		event.Start,
		event.End,
		event.Reoccurrence,
		event.ForUnloggedins,
		event.ForUnverifieds,
		event.InviteOnly,
		event.HostId,
		event.AddressId,
	)
	if err != nil {
		return err
	}

	return nil
}
