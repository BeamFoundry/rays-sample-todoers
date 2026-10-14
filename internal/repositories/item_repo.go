package repositories

import (
	. "github.com/BeamFoundry/rays-orm/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays/pkg/lang"
)

type ItemRepositoryImpl struct {
	OrmRepository[Item, uint]
	findAllByTodoIDOrderByID                 func(uint) []Item
	findAllByTodoIDOrderByRank               func(uint) []Item
	findAllByTodoIDAndRankGreaterThanOrEqual func(uint, uint) []Item
	findOneByTodoIDAndItemID                 func(uint, uint) Optional[*Item]
	findOneByTodoIDAndRank                   func(uint, uint) Optional[*Item]
}

var _ ItemRepository = &ItemRepositoryImpl{}

func (this *ItemRepositoryImpl) FindAllByTodoID(todoID uint) []Item {
	return this.findAllByTodoIDOrderByID(todoID)
}

func (this *ItemRepositoryImpl) FindAllByTodoIDOrderByRank(todoID uint) []Item {
	return this.findAllByTodoIDOrderByRank(todoID)
}

func (this *ItemRepositoryImpl) FindAllByTodoIDAndRankGreaterThanOrEqual(todoID, rank uint) []Item {
	return this.findAllByTodoIDAndRankGreaterThanOrEqual(todoID, rank)
}

func (this *ItemRepositoryImpl) FindOneByTodoIDAndItemID(todoID, itemID uint) Optional[*Item] {
	return this.findOneByTodoIDAndItemID(todoID, itemID)
}

func (this *ItemRepositoryImpl) FindOneByTodoIDAndRank(todoID, rank uint) Optional[*Item] {
	return this.findOneByTodoIDAndRank(todoID, rank)
}
