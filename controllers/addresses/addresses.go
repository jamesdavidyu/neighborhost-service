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
