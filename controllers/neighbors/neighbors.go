package neighbors

import (
	"database/sql"
	"fmt"

	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetNeighborWithEmail(email string) (*types.Neighbors, error) {
	var neighbor types.Neighbors

	rows, err := s.db.Query(
		`SELECT * FROM neighbors
		WHERE email = $1`,
		email,
	)
	// need to handle errors better
	if err != nil {
		return nil, err
	}
	// defer rows.Close()

	// if !rows.Next() {
	// 	return nil, sql.ErrNoRows
	// }

	if err := rows.Scan(
		&neighbor.Id,
		&neighbor.Email,
		&neighbor.Username,
		&neighbor.Zipcode,
		&neighbor.Password,
		&neighbor.Verified,
		&neighbor.NeighborhoodId,
		&neighbor.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &neighbor, nil
}

func (s *Store) GetNeighborWithUsername(username string) (*types.Neighbors, error) {
	var neighbor types.Neighbors

	rows, err := s.db.Query(
		`SELECT * FROM neighbors
		WHERE username = $1`,
		username,
	)
	if err != nil {
		return nil, err
	}
	// defer rows.Close()

	// if !rows.Next() {
	// 	return nil, sql.ErrNoRows
	// }

	if err := rows.Scan(
		&neighbor.Id,
		&neighbor.Email,
		&neighbor.Username,
		&neighbor.Zipcode,
		&neighbor.Password,
		&neighbor.Verified,
		&neighbor.NeighborhoodId,
		&neighbor.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &neighbor, nil
}

func (s *Store) GetNeighborWithEmailOrUsername(emailOrUsername string) (*types.Neighbors, error) {
	rows, err := s.db.Query(
		`SELECT * FROM neighbors
		WHERE email = $1 OR username = $1`,
		emailOrUsername,
	)
	if err != nil {
		return nil, err
	}

	neighbor := new(types.Neighbors)
	for rows.Next() {
		neighbor, err = utils.ScanRowIntoNeighbor(rows)
		if err != nil {
			return nil, err
		}
	}

	if neighbor.Id == 0 {
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
		neighbor, err = utils.ScanRowIntoNeighbor(rows)
		if err != nil {
			return nil, err
		}
	}

	if neighbor.Id == 0 {
		return nil, fmt.Errorf("error")
	}

	return neighbor, nil
}

func (s *Store) UpdateZipcodeWithId(neighbor types.Neighbors) error {
	_, err := s.db.Exec(
		`UPDATE neighbors
		SET zipcode = $1
		WHERE id = $2`,
		neighbor.Zipcode, neighbor.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdatePasswordWithId(neighbor types.Neighbors) error {
	_, err := s.db.Exec(
		`UPDATE neighbors
		SET password = $1
		WHERE id = $2`,
		neighbor.Password, neighbor.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateNeighbor(neighbor types.Neighbors) error {
	_, err := s.db.Exec(
		`INSERT INTO neighbors (email, username, zipcode, password)
		VALUES ($1, $2, $3, $4)`,
		neighbor.Email,
		neighbor.Username,
		neighbor.Zipcode,
		neighbor.Password,
	)
	if err != nil {
		return err
	}

	return nil
}
