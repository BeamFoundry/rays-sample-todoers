package domains

import (
	"time"
)

//nolint:govet // model
type Item struct {
	ID        uint
	Text      string
	Done      bool
	TodoID    uint
	ItemID    uint
	Rank      uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
