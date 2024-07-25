package events

import (
	"database/sql"
	"time"
)

type EventAddresses struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Start            time.Time `json:"start"`
	End              time.Time `json:"end"`
	Reoccurrence     string    `json:"reoccurrence"`
	ForUnloggedins   bool      `json:"forUnloggedins"`
	ForUnverifieds   bool      `json:"forUnverifieds"`
	InviteOnly       bool      `json:"inviteOnly"`
	HostId           int       `json:"hostId"`
	AddressId        int       `json:"addressId"`
	CreatedAt        time.Time `json:"createdAt"`
	AddressAddressId int       `json:"addressAddressId"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	Address          string    `json:"address"`
	City             string    `json:"city"`
	State            string    `json:"state"`
	Zipcode          string    `json:"zipcode"`
	NeighborId       int       `json:"neighborId"`
	NeighborhoodId   int       `json:"neighborhoodId"`
	RecordedAt       time.Time `json:"recordedAt"`
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetEventsByZipcode(zipcode string) (*EventAddresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM events e
		LEFT OUTER JOIN addresses a ON a.id = e.address_id
		WHERE a.zipcode = $1`, zipcode,
	)
	if err != nil {
		return nil, err
	}

	events := new(EventAddresses)
	for rows.Next() {
		events, err = scanRowIntoEvents(rows)
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}

func scanRowIntoEvents(rows *sql.Rows) (*EventAddresses, error) {
	events := new(EventAddresses)

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
