package core

import (
	. "github.com/BeamFoundry/rays-orm/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays/pkg/lang"
)

type ItemRepository interface {
	OrmRepositoryInterface[Item, uint]
	FindAllByTodoID(todoId uint) []Item
	FindAllByTodoIDOrderByRank(todoId uint) []Item
	FindAllByTodoIDAndRankGreaterThanOrEqual(uint, uint) []Item
	FindOneByTodoIDAndItemID(uint, uint) Optional[*Item]
	FindOneByTodoIDAndRank(uint, uint) Optional[*Item]
}

type TodoRepository interface {
	OrmRepositoryInterface[Todo, uint]
	FindAllByUserID(uint) []Todo
	FindAllByUserIDOrderByRank(uint) []Todo
	FindAllByUserIDWithItems(uint) []Todo
	FindOneWithItemsByID(uint) Optional[*Todo]
	FindOneByUserIDAndTodoID(uint, uint) Optional[*Todo]
}

type UserRepository interface {
	OrmRepositoryInterface[User, uint]
	FindByUsername(string) Optional[*User]
	FindAllUsersByUsername(string, ...OrmQueryFunc) []User
}
