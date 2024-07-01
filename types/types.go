package types

import "time"

type Neighborhoods struct {
	ID           int       `json:"id"`
	Neighborhood string    `json:"neighborhood"`
	CreatedAt    time.Time `json:"createdAt"`
}
