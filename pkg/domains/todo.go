package domains

import (
	"time"
)

//nolint:govet // model
type Todo struct {
	ID        uint
	UserID    uint
	Title     string
	TodoID    uint
	Rank      uint
	Items     []Item
	CreatedAt time.Time
	UpdatedAt time.Time
}
