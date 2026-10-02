package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays-orm"
	. "github.com/BeamFoundry/rays/pkg/lang"
)

type TodoRepository interface {
	OrmRepositoryInterface[Todo, uint]
	FindAllByUserID(uint) []Todo
	FindAllByUserIDOrderByRank(uint) []Todo
	FindAllByUserIDWithItems(uint) []Todo
	FindOneWithItemsByID(uint) Optional[*Todo]
	FindOneByUserIDAndTodoID(uint, uint) Optional[*Todo]
}

type TodoRepositoryImpl struct {
	OrmRepositoryComposite[Todo, uint]
	findOneByUserIDAndTodoID   func(uint, uint) Optional[*Todo]
	findOneByID                func(uint, ...OrmQueryFunc) Optional[*Todo]
	findAllByUserIDOrderByID   func(uint) []Todo
	findAllByUserIDOrderByRank func(uint, ...OrmQueryFunc) []Todo
	//+ findAllByUserIDOrderByRankWithItemsOrderByRank func(uint) []Todo
}

var _ TodoRepository = &TodoRepositoryImpl{}

var withItems = func(db *gorm.DB) *gorm.DB {
	return db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		columns := []clause.OrderByColumn{
			{
				Column: clause.Column{Name: "items.rank"},
			},
		}
		return db.Order(clause.OrderBy{Columns: columns})
	})
}

func (this *TodoRepositoryImpl) FindAllByUserID(userID uint) []Todo {
	return this.findAllByUserIDOrderByID(userID)
}

func (this *TodoRepositoryImpl) FindAllByUserIDOrderByRank(userID uint) []Todo {
	return this.findAllByUserIDOrderByRank(userID)
}

func (this *TodoRepositoryImpl) FindAllByUserIDWithItems(userID uint) []Todo {
	return this.findAllByUserIDOrderByRank(userID, withItems)
}

func (this *TodoRepositoryImpl) FindOneWithItemsByID(todoId uint) Optional[*Todo] {
	return this.findOneByID(todoId, withItems)
}

func (this *TodoRepositoryImpl) FindOneByUserIDAndTodoID(userID, todoID uint) Optional[*Todo] {
	return this.findOneByUserIDAndTodoID(userID, todoID)
}
