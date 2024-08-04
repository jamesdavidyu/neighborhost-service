package friends

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

func (s *Store) GetFriendsByNeighborId(neighborId int) ([]types.Friends, error) {
	rows, err := s.db.Query(
		`SELECT * FROM friends
		WHERE neighbor_id = $1`, neighborId,
	)
	if err != nil {
		return nil, err
	}

	friends := make([]types.Friends, 0)
	for rows.Next() {
		friend, err := utils.ScanRowIntoFriends(rows)
		if err != nil {
			return nil, err
		}
		friends = append(friends, *friend)
	}

	return friends, nil
}

func (s *Store) CreateFriendRequest(friendRequest types.FriendRequests) error {
	_, err := s.db.Exec(
		`INSERT INTO friend_requests (
			neighbor_id,
			requesting_friend_id,
			status
		)
		VALUES ($1, $2, $3)`,
		friendRequest.NeighborId,
		friendRequest.RequestingFriendId,
		friendRequest.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetFriendRequestsByNeighborId(neighborId int) ([]types.FriendRequests, error) {
	rows, err := s.db.Query(
		`SELECT * FROM friend_requests
		WHERE neighbor_id = $1`, neighborId,
	)
	if err != nil {
		return nil, err
	}

	friendRequests := make([]types.FriendRequests, 0)
	for rows.Next() {
		friendRequest, err := utils.ScanRowIntoFriendRequests(rows)
		if err != nil {
			return nil, err
		}
		friendRequests = append(friendRequests, *friendRequest)
	}

	return friendRequests, nil
}
