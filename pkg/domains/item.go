package domains

import (
	"time"
)

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

func (this Item) Id() uint {
	return this.ID
}
