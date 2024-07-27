package addresses

import (
	"database/sql"

	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/utils"
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

func (s *Store) GetAddressIdByAddress(
	firstName string,
	lastName string,
	address string,
	city string,
	state string,
	zipcode string,
	neighborId int,
) (*types.Addresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM addresses
		WHERE first_name = $1
		AND last_name = $2
		AND address = $3
		AND city = $4
		AND state = $5
		AND zipcode = $6
		AND neighbor_id = $7`,
		firstName,
		lastName,
		address,
		city,
		state,
		zipcode,
		neighborId,
	)
	if err != nil {
		return nil, err
	}

	// if !rows.Next() {
	// 	return nil, sql.ErrNoRows
	// }

	addresses := new(types.Addresses)
	for rows.Next() {
		addresses, err = utils.ScanRowIntoAddresses(rows)
		if err != nil {
			return nil, err
		}
	}

	return addresses, nil
}

func (s *Store) GetAddressByNeighborId(id int) (*types.Addresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM addresses
		WHERE neighbor_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}

	addresses := new(types.Addresses)
	for rows.Next() {
		addresses, err = utils.ScanRowIntoAddresses(rows)
		if err != nil {
			return nil, err
		}
	}

	return addresses, nil
}

// needed??
func (s *Store) GetAddressesByZipcode(zipcode string) (*types.Addresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM addresses
		WHERE zipcode = $1`, zipcode,
	)
	if err != nil {
		return nil, err
	}

	addresses := new(types.Addresses) // change to make like in neighborhoods - do this for any get requests for more than one row
	for rows.Next() {
		addresses, err = utils.ScanRowIntoAddresses(rows)
		if err != nil {
			return nil, err
		}
	}

	return addresses, nil
}
