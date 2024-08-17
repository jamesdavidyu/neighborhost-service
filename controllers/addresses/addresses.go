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
			type,
			neighbor_id,
			neighborhood_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		address.FirstName,
		address.LastName,
		address.Address,
		address.City,
		address.State,
		address.Zipcode,
		address.Type,
		address.NeighborId,
		address.NeighborhoodId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAddressIdByAddress(
	address string,
	city string,
	state string,
	zipcode string,
	Type string,
	neighborId int,
) (*types.Addresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM addresses
		WHERE address = $1
		AND city = $2
		AND state = $3
		AND zipcode = $4
		AND type = $5
		AND neighbor_id = $6`,
		address,
		city,
		state,
		zipcode,
		Type,
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
		WHERE neighbor_id = $1
		AND type = 'Home'`, id,
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

func (s *Store) GetAddressesByNeighborId(id int) ([]types.Addresses, error) {
	rows, err := s.db.Query(
		`SELECT * FROM addresses
		WHERE neighbor_id = $1
		ORDER BY address`, id,
	)
	if err != nil {
		return nil, err
	}

	addresses := make([]types.Addresses, 0)
	for rows.Next() {
		address, err := utils.ScanRowIntoAddresses(rows)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, *address)
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
