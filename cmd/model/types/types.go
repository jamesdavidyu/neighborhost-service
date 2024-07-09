package types

import "time"

type NeighborStore interface {
	GetNeighborWithEmailOrUsername(emailOrUsername string) (*Neighbors, error)
	GetNeighborById(id int) (*Neighbors, error)
	CreateNeighbor(Neighbors) error
	GetNeighborWithEmail(email string) (*Neighbors, error)
	GetNeighborWithUsername(username string) (*Neighbors, error)
	// UpdatePassword(id int) (*Neighbors, error)
	// UpdateZipcode(id int) (*Neighbors, error)
}

type EventStore interface {
	GetEventsByZipcode(zipcode string) (*Events, error)
}

type AddressStore interface {
	CreateAddress(Addresses) error
}

type NeighborhoodStore interface {
	GetNeighborhoods() ([]Neighborhoods, error)
	CreateNeighborhood(Neighborhoods) error
}

type Neighborhoods struct {
	Id           int       `json:"id"`
	Neighborhood string    `json:"neighborhood"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Neighbors struct {
	Id             int       `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Zipcode        string    `json:"zipcode"`
	Password       string    `json:"-"`
	Verified       bool      `json:"verified"`
	NeighborhoodId int       `json:"neighborhoodId"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Register struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Zipcode  string `json:"zipcode" validate:"required,min=5,max=5"`
	Password string `json:"password" validate:"required,min=8"`
}

type Login struct {
	EmailOrUsername string `json:"emailOrUsername" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
}

type Addresses struct {
	Id             int       `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	State          string    `json:"state"`
	Zipcode        string    `json:"zipcode"`
	NeighborId     int       `json:"neighborId"`
	NeighborhoodId int       `json:"neighborhoodId"`
	RecordedAt     time.Time `json:"recordedAt"`
}

type AddressPayload struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Address        string `json:"address"`
	City           string `json:"city"`
	State          string `json:"state"`
	Zipcode        string `json:"zipcode"`
	NeighborId     int    `json:"neighborId"`
	NeighborhoodId int    `json:"neighborhoodId"`
}

type Profiles struct {
	Id         int `json:"id"`
	NeighborId int `json:"neighborId"`
	// demographic fields... bio, race, age
}

type Events struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
	Reoccurrence   string    `json:"reoccurrence"`
	ForUnloggedins bool      `json:"forUnloggedins"`
	ForUnverifieds bool      `json:"forUnverifieds"`
	InviteOnly     bool      `json:"inviteOnly"`
	HostId         int       `json:"hostId"`
	AddressId      int       `json:"addressId"`
	CreatedAt      time.Time `json:"createdAt"`
}

type EventInvites struct {
	Id                int       `json:"id"`
	EventId           int       `json:"eventId"`
	InvitedNeighborId int       `json:"invitedNeighborId"`
	InvitedAt         time.Time `json:"invitedAt"`
}

type Friends struct {
	Id                int       `json:"id"`
	NeighborId        int       `json:"neighborId"`
	NeighborsFriendId int       `json:"neighborsFriendId"`
	FriendedAt        time.Time `json:"friendedAt"`
}

type FriendRequests struct {
	Id                 int       `json:"id"`
	NeighborId         int       `json:"neighborId"`
	RequestingFriendId int       `json:"requestingFriendId"`
	FriendRequestedAt  time.Time `json:"friendRequestedAt"`
}
