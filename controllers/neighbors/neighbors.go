package neighbors

import (
	"database/sql"
	"fmt"

	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetNeighborWithEmailOrUsername(emailOrUsername string) (*types.Neighbors, error) {
	rows, err := s.db.Query(`SELECT * FROM neighbors
							WHERE email = $1 OR username = $1`, emailOrUsername)
	if err != nil {
		return nil, err
	}

	neighbor := new(types.Neighbors)
	for rows.Next() {
		neighbor, err = scanRowIntoNeighbor(rows)
		if err != nil {
			return nil, err
		}
	}

	if neighbor.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return neighbor, nil
}

func (s *Store) GetNeighborById(id int) (*types.Neighbors, error) {
	rows, err := s.db.Query("SELECT * FROM neighbors WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	neighbor := new(types.Neighbors)
	for rows.Next() {
		neighbor, err = scanRowIntoNeighbor(rows)
		if err != nil {
			return nil, err
		}
	}

	if neighbor.ID == 0 {
		return nil, fmt.Errorf("error")
	}

	return neighbor, nil
}

func (s *Store) CreateNeighbor(neighbor types.Neighbors) error {
	_, err := s.db.Exec(`INSERT INTO neighbors (email, username, zipcode, password)
							VALUES ($1, $2, $3, $4)`, neighbor.Email, neighbor.Username, neighbor.Zipcode, neighbor.Password)
	if err != nil {
		return err
	}

	return nil
}

func scanRowIntoNeighbor(rows *sql.Rows) (*types.Neighbors, error) {
	neighbor := new(types.Neighbors)

	err := rows.Scan(&neighbor.ID, &neighbor.Email, &neighbor.Username, &neighbor.Zipcode, &neighbor.Password, &neighbor.Verified, &neighbor.NeighborhoodID, &neighbor.CreatedAt)
	if err != nil {
		return nil, err
	}

	return neighbor, nil
}
