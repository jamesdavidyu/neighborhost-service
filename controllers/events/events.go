package events

import (
	"database/sql"

	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetPublicEvents() ([]types.Events, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events
		WHERE for_unloggedins = TRUE`,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.Events, 0)
	for rows.Next() {
		event, err := scanRowIntoPublicEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (s *Store) GetEventsByZipcode(zipcode string) ([]types.EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1`, zipcode,
	)
	if err != nil {
		return nil, err
	}

	events := make([]types.EventAddresses, 0)
	for rows.Next() {
		event, err := scanRowIntoNeighborEvents(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func scanRowIntoPublicEvents(rows *sql.Rows) (*types.Events, error) {
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

func scanRowIntoNeighborEvents(rows *sql.Rows) (*types.EventAddresses, error) {
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
