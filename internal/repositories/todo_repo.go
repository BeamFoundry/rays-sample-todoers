package repositories

import (
	. "github.com/BeamFoundry/rays-orm/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays/pkg/lang"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TodoRepositoryImpl struct {
	OrmRepository[Todo, uint]
	findOneByUserIDAndTodoID   func(uint, uint) Optional[*Todo]
	findOneByID                func(uint, ...OrmQueryFunc) Optional[*Todo]
	findAllByUserIDOrderByID   func(uint) []Todo
	findAllByUserIDOrderByRank func(uint, ...OrmQueryFunc) []Todo
	//+ findAllByUserIDOrderByRankWithItemsOrderByRank func(uint) []Todo
}

var _ TodoRepository = &TodoRepositoryImpl{}

func (this *TodoRepositoryImpl) FindAllByUserID(userID uint) []Todo {
	return this.findAllByUserIDOrderByID(userID)
}

func (this *TodoRepositoryImpl) FindAllByUserIDOrderByRank(userID uint) []Todo {
	return this.findAllByUserIDOrderByRank(userID)
}

func (this *TodoRepositoryImpl) FindAllByUserIDWithItems(userID uint) []Todo {
	return this.findAllByUserIDOrderByRank(userID, loadItems)
}

func (this *TodoRepositoryImpl) FindOneWithItemsByID(todoId uint) Optional[*Todo] {
	return this.findOneByID(todoId, loadItems)
}

func (this *TodoRepositoryImpl) FindOneByUserIDAndTodoID(userID, todoID uint) Optional[*Todo] {
	return this.findOneByUserIDAndTodoID(userID, todoID)
}

func loadItems(db *gorm.DB) *gorm.DB {
	return db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		columns := []clause.OrderByColumn{
			{
				Column: clause.Column{Name: "items.rank"},
			},
		}
		return db.Order(clause.OrderBy{Columns: columns})
	})
}
