package domains

import (
	"time"
	//"gorm.io/gorm"
	//"gorm.io/gorm/clause"
)

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

func (this Todo) Id() uint {
	return this.ID
}

/*
func (this *Todo) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("todo_id = ?", this.ID).Delete(&Item{})
	return
}
*/
