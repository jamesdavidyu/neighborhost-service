/*
1. PUBLIC
2. ZIPCODE
3. NEIGHBORHOOD
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
		AND start >= $1`, time.Now(),
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

func (s *Store) GetEventsByZipcode(zipcode string, start time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1
		AND start >= $2`, zipcode, start,
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

func (s *Store) ZipcodeEventsOnDate(zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1
		AND start = $2`, zipcode, dateTime,
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

func (s *Store) ZipcodeEventsBeforeDate(zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1
		AND start < $2`, zipcode, dateTime,
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

func (s *Store) ZipcodeEventsAfterDate(zipcode string, dateTime time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1
		AND start >= $2`, zipcode, dateTime,
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

/* 3. NEIGHBORHOOD */

func (s *Store) GetEventsByNeighborhoodId(id int, start time.Time) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.neighborhood_id = $1
		AND start >= $2`, id, start,
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
