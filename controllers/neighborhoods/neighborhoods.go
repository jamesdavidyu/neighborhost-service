package neighborhoods

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
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
		neighborhood, err := scanRowsIntoNeighborhood(rows)
		if err != nil {
			return nil, err
		}
		neighborhoods = append(neighborhoods, *neighborhood)
	}

	return neighborhoods, nil
}

func CreateNeighborhood(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var neighborhood types.Neighborhoods
		json.NewDecoder(r.Body).Decode(&neighborhood)

		_, err := db.Exec(`INSERT INTO neighborhoods (neighborhood)
							VALUES ($1)`, neighborhood.Neighborhood)
		if err != nil {
			fmt.Println(err)
		}

		json.NewEncoder(w).Encode(neighborhood)
	}
}

func scanRowsIntoNeighborhood(rows *sql.Rows) (*types.Neighborhoods, error) {
	neighborhood := new(types.Neighborhoods)

	err := rows.Scan(
		&neighborhood.ID,
		&neighborhood.Neighborhood,
		&neighborhood.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return neighborhood, nil
}
