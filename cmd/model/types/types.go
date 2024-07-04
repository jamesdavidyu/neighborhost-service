package types

import "time"

type NeighborStore interface {
	GetNeighborWithEmailOrUsername(emailOrUsername string) (*Neighbors, error)
	GetNeighborById(id int) (*Neighbors, error)
	CreateNeighbor(Neighbors) error
}

type NeighborhoodStore interface {
	GetNeighborhoods() ([]Neighborhoods, error)
}

type Neighborhoods struct {
	ID           int       `json:"id"`
	Neighborhood string    `json:"neighborhood"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Neighbors struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Zipcode        string    `json:"zipcode"`
	Password       string    `json:"-"`
	Verified       bool      `json:"verified"`
	NeighborhoodID int       `json:"neighborhoodId"`
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
