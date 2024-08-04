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
			requested_friend_id,
			status
		)
		VALUES ($1, $2, $3)`,
		friendRequest.NeighborId,
		friendRequest.RequestedFriendId,
		friendRequest.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetFriendRequestsByNeighborId(requestedFriendId int) ([]types.PendingFriendRequests, error) {
	rows, err := s.db.Query(
		`SELECT * FROM friend_requests f
		JOIN neighbors n ON n.id = f.requested_friend_id
		WHERE f.requested_friend_id = $1
		AND f.status = 'pending'`, requestedFriendId,
	)
	if err != nil {
		return nil, err
	}

	friendRequests := make([]types.PendingFriendRequests, 0)
	for rows.Next() {
		friendRequest, err := utils.ScanRowIntoFriendRequests(rows)
		if err != nil {
			return nil, err
		}
		friendRequests = append(friendRequests, *friendRequest)
	}

	return friendRequests, nil
}

// respond to friend request controller
func (s *Store) UpdateFriendRequest(friendRequest types.FriendRequests) error {
	_, err := s.db.Exec(
		`UPDATE friend_requests
		SET status = $1
		WHERE neighbor_id = $2
		AND requested_friend_id = $3`,
		friendRequest.Status,
		friendRequest.NeighborId,
		friendRequest.RequestedFriendId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateFriend(friends types.Friends) error {
	_, err := s.db.Exec(
		`INSERT INTO friends (
			neighbor_id,
			neighbors_friend_id
		)
		VALUES ($1, $2)`, friends.NeighborId, friends.NeighborsFriendId,
	)
	if err != nil {
		return err
	}

	return nil
}
