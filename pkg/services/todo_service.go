package services

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"gorm.io/gorm"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/i18n"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-orm"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/repositories"
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

type TodoServiceImpl struct {
	Service
	transactionManager *GormTxManager `@:"Inject"`
	todoRepository     TodoRepository `@:"Inject"`
	itemRepository     ItemRepository `@:"Inject"`

	localeContextHolder LocaleContextHolder `@:"Inject"`
	messageSource       MessageSource       `@:"Inject"`
}

var _ TodoService = &TodoServiceImpl{}

func (this *TodoServiceImpl) CreateTodolist(user *User) (*Todo, error) {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()

	var ids []uint
	var id uint = 1

	this.todoRepository.Query(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", user.Id())
	}).Pluck("todo_id", &ids)

	if len(ids) > 0 {
		id = slices.Max(ids) + 1
	}

	item_title := this.TranslateWithDefault("item.title", "Item 1", "1")
	todo_title := this.TranslateWithDefault("todo.title", fmt.Sprintf("Todo %d", id), strconv.FormatUint(uint64(id), 10))
	item := &Item{ItemID: 1, Text: item_title, Rank: 1}
	todo := &Todo{TodoID: id, Title: todo_title, Items: []Item{*item}, UserID: user.Id(), Rank: uint(len(ids))}
	if t, err := this.todoRepository.Save(todo); err != nil {
		tx.Rollback()
		return nil, err
	} else {
		todo = t
	}

	tx.Commit()
	return todo, nil
}

func (this *TodoServiceImpl) TranslateWithDefault(message, def, item string) string {
	locale := this.localeContextHolder.Get()
	bundle := this.messageSource.FindMessageBundle(locale)
	msg, err := bundle.Translate(message, item)
	if err != nil {
		return def
	}
	return msg
}

func (this *TodoServiceImpl) GetUserTodoLists(user *User) []Todo {
	return this.todoRepository.FindAllByUserIDWithItems(user.Id())
}

func (this *TodoServiceImpl) GetUserList(user *User, todoID uint) Optional[*Todo] {
	return this.todoRepository.FindOneByUserIDAndTodoID(user.Id(), todoID)
}

func (this *TodoServiceImpl) GetListWithItems(id uint) Optional[*Todo] {
	return this.todoRepository.FindOneWithItemsByID(id)
}

func (this *TodoServiceImpl) SaveList(todo *Todo) (*Todo, error) {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()
	if t, err := this.todoRepository.Save(todo); err != nil {
		tx.Rollback()
		return nil, err
	} else {
		tx.Commit()
		return t, nil
	}
}

func (this *TodoServiceImpl) DeleteList(todo *Todo) error {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()
	if err := this.todoRepository.DeleteById(todo.Id()); err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}

func (this *TodoServiceImpl) ReorderLists(user *User, order []uint) error {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()
	todos := this.todoRepository.FindAllByUserID(user.Id())
	if len(todos) != len(order) {
		tx.Rollback()
		return NewError("InvalidListCount", "Current list count differs from argument")
	}
	for index, todo := range todos {
		if err := tx.Model(&todo).Update("rank", order[index]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (this *TodoServiceImpl) ReorderListItems(todo *Todo, order []uint) error {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()
	items := this.itemRepository.FindAllByTodoID(todo.Id())
	if len(items) != len(order) {
		tx.Rollback()
		return NewError("InvalidListCount", "Current list count differs from argument")
	}
	for index, item := range items {
		if err := tx.Model(&item).Update("rank", order[index]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (this *TodoServiceImpl) MoveListItem(todo *Todo, oldRank, to, newRank uint) error {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()
	if item, ok := this.itemRepository.FindOneByTodoIDAndRank(todo.Id(), oldRank).Get(); !ok {
		return NewError("itemRepositoryError", "Unable to find source item")
	} else {
		items := this.itemRepository.FindAllByTodoIDAndRankGreaterThanOrEqual(to, newRank)
		for _, i := range items {
			i.Rank++
			if _, err := this.itemRepository.Save(&i); err != nil {
				tx.Rollback()
				return err
			}
		}
		item.TodoID = to
		item.Rank = newRank
		if _, err := this.itemRepository.Save(item); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}
}

func (this *TodoServiceImpl) GetTodoItem(todo *Todo, itemID uint) Optional[*Item] {
	return this.itemRepository.FindOneByTodoIDAndItemID(todo.Id(), itemID)
}

func (this *TodoServiceImpl) GetItem(id uint) Optional[*Item] {
	return this.itemRepository.FindById(id)
}

func (this *TodoServiceImpl) GetItems(todo *Todo) []Item {
	return this.itemRepository.FindAllByTodoIDOrderByRank(todo.Id())
}

func (this *TodoServiceImpl) AddItem(todo *Todo) (*Item, error) {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()
	var ids []uint
	var id uint = 1

	if todo == nil {
		err := NewError("TodoServiceError", "Todo is empty")
		err.LogError()
		return nil, err
	}

	this.itemRepository.Query(func(db *gorm.DB) *gorm.DB {
		return db.Where("todo_id = ?", todo.Id())
	}).Pluck("item_id", &ids)

	if len(ids) > 0 {
		id = slices.Max(ids) + 1
	}

	item_title := this.TranslateWithDefault("item.title", fmt.Sprintf("Item %d", id), strconv.FormatUint(uint64(id), 10))
	item := &Item{TodoID: todo.Id(), ItemID: id, Text: item_title, Rank: uint(len(ids) + 1)}
	if i, err := this.itemRepository.Save(item); err != nil {
		tx.Rollback()
		return nil, err
	} else {
		tx.Commit()
		return i, nil
	}
}

func (this *TodoServiceImpl) SaveItem(item *Item) (*Item, error) {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()
	if t, err := this.itemRepository.Save(item); err != nil {
		tx.Rollback()
		return nil, err
	} else {
		tx.Commit()
		return t, nil
	}
}

func (this *TodoServiceImpl) DeleteItem(item *Item) error {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()
	if err := this.itemRepository.DeleteById(item.Id()); err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}
