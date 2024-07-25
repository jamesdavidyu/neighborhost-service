package addresses

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

func (s *Store) CreateAddress(address types.Addresses) error {
	_, err := s.db.Exec(
		`INSERT INTO addresses (
			first_name,
			last_name,
			address,
			city, 
			state,
			zipcode,
			neighbor_id,
			neighborhood_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		address.FirstName,
		address.LastName,
		address.Address,
		address.City,
		address.State,
		address.Zipcode,
		address.NeighborId,
		address.NeighborhoodId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAddressesByZipcode(zipcode string) (*types.Addresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM addresses
		WHERE zipcode = $1`, zipcode,
	)
	if err != nil {
		return nil, err
	}

	addresses := new(types.Addresses) // change to make like in neighborhoods - do this for any get requests
	for rows.Next() {
		addresses, err = scanRowIntoAddresses(rows)
		if err != nil {
			return nil, err
		}
	}

	return addresses, nil
}

func scanRowIntoAddresses(rows *sql.Rows) (*types.Addresses, error) {
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
