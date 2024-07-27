package zipcodes

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

func (s *Store) GetZipcodeData(zipcode string) (*types.Zipcodes, error) {
	rows, err := s.db.Query(
		`SELECT * FROM zipcodes
		WHERE zipcode = $1`, zipcode,
	)
	if err != nil {
		return nil, err
	}

	zipcodeData := new(types.Zipcodes)
	for rows.Next() {
		zipcodeData, err = utils.ScanRowIntoZipcodes(rows)
		if err != nil {
			return nil, err
		}
	}

	return zipcodeData, nil
}

func (s *Store) ValidateZipcode(city string, state string, zipcode string) (*types.Zipcodes, error) {
	rows, err := s.db.Query(
		`SELECT * FROM zipcodes
		WHERE city = $1
		AND state = $2
		AND zipcode = $3`, city, state, zipcode,
	)
	if err != nil {
		return nil, err
	}

	// if !rows.Next() {
	// 	return nil, sql.ErrNoRows
	// }

	zipcodes := new(types.Zipcodes)
	for rows.Next() {
		zipcodes, err = utils.ScanRowIntoZipcodes(rows)
		if err != nil {
			return nil, err
		}
	}

	return zipcodes, nil
}
