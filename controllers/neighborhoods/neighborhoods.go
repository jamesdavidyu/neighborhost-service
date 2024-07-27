package neighborhoods

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

func (s *Store) GetNeighborhoods() ([]types.Neighborhoods, error) {
	rows, err := s.db.Query("SELECT * FROM neighborhoods")
	if err != nil {
		return nil, err
	}

	neighborhoods := make([]types.Neighborhoods, 0)
	for rows.Next() {
		neighborhood, err := utils.ScanRowsIntoNeighborhood(rows)
		if err != nil {
			return nil, err
		}
		neighborhoods = append(neighborhoods, *neighborhood)
	}

	return neighborhoods, nil
}

func (s *Store) CreateNeighborhood(neighborhood types.Neighborhoods) error {
	_, err := s.db.Exec(
		`INSERT INTO neighborhoods (neighborhood)
		VALUES ($1)`,
		neighborhood.Neighborhood,
	)
	if err != nil {
		return err
	}

	return nil
}
