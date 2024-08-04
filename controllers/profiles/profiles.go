package profiles

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

func (s *Store) GetProfileByNeighborId(neighborId int) (*types.Profiles, error) {
	rows, err := s.db.Query(
		`SELECT 
			b.neighbor_id,
			b.bio,
			dob.date_of_birth,
			dob.date_of_birth_public,
			g.gender,
			g.gender_public,
			r.race,
			r.race_public,
			e.ethnicity,
			e.ethnicity_public,
			rs.relationship_status,
			rs.relationship_status_public,
			re.religion,
			re.religion_public,
			p.politics,
			p.politics_public
		FROM bios b
		JOIN dates_of_birth dob ON dob.neighbor_id = b.neighbor_id
		JOIN genders g ON g.neighbor_id = b.neighbor_id
		JOIN races r ON r.neighbor_id = b.neighbor_id
		JOIN ethnicities e ON e.neighbor_id = b.neighbor_id
		JOIN relationship_statuses rs ON rs.neighbor_id = b.neighbor_id
		JOIN religions re ON re.neighbor_id = b.neighbor_id
		JOIN politics p On p.neighbor_id = b.neighbor_id
		WHERE b.neighbor_id = $1`, neighborId,
	)
	if err != nil {
		return nil, err
	}

	profile := new(types.Profiles)
	for rows.Next() {
		profile, err = utils.ScanRowIntoProfiles(rows)
		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}
