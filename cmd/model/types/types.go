package types

import "time"

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
	Zipcode  string `json:"zipcode" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
