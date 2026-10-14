package core

import (
	. "github.com/BeamFoundry/rays-web-security/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
)

type TodoService interface {
	CreateTodolist(*User) (*Todo, error)
	GetUserTodoLists(*User) []Todo
	GetUserList(*User, uint) Optional[*Todo]
	GetListWithItems(uint) Optional[*Todo]
	SaveList(*Todo) (*Todo, error)
	DeleteList(*Todo) error
	ReorderLists(user *User, order []uint) error
	ReorderListItems(todo *Todo, item []uint) error
	GetTodoItem(*Todo, uint) Optional[*Item]
	GetItem(uint) Optional[*Item]
	GetItems(*Todo) []Item
	AddItem(*Todo) (*Item, error)
	SaveItem(*Item) (*Item, error)
	DeleteItem(*Item) error
	MoveListItem(todo *Todo, oldRank, to, newRank uint) error
	// GetUserListWithItems(*User, uint) Optional[*Todo]
	// GetList(uint) Optional[*Todo]
}

type UserService interface {
	UserDetailsService
	AddUser(name, username, password string) (*User, error)
	GetFullUser(id uint) Optional[*User]
	UpdateName(name string, user *User) (*User, error)
	UpdatePassword(password string, user *User) (*User, error)
	FindAllUsersByUsername(query string, page, pageSize int) []User
	CountAll() int64
}
