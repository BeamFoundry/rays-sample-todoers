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

	. "github.com/BeamFoundry/rays-orm/pkg/core"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
)

type TodoServiceImpl struct {
	Service
	transactional  GormTransactional `@:"Inject"`
	todoRepository TodoRepository    `@:"Inject"`
	itemRepository ItemRepository    `@:"Inject"`

	localeContextHolder LocaleContextHolder `@:"Inject"`
	messageSource       MessageSource       `@:"Inject"`
}

var _ TodoService = &TodoServiceImpl{}

func (this *TodoServiceImpl) CreateTodolist(user *User) (*Todo, error) {
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()

	var ids []uint
	var id uint = 1

	this.todoRepository.Query(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", user.ID)
	}).Pluck("todo_id", &ids)

	if len(ids) > 0 {
		id = slices.Max(ids) + 1
	}

	item_title := this.TranslateWithDefault("item.title", "Item 1", "1")
	todo_title := this.TranslateWithDefault("todo.title", fmt.Sprintf("Todo %d", id), strconv.FormatUint(uint64(id), 10))
	item := &Item{ItemID: 1, Text: item_title, Rank: 1}
	todo := &Todo{TodoID: id, Title: todo_title, Items: []Item{*item}, UserID: user.ID, Rank: uint(len(ids))}
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
	return this.todoRepository.FindAllByUserIDWithItems(user.ID)
}

func (this *TodoServiceImpl) GetUserList(user *User, todoID uint) Optional[*Todo] {
	return this.todoRepository.FindOneByUserIDAndTodoID(user.ID, todoID)
}

func (this *TodoServiceImpl) GetListWithItems(id uint) Optional[*Todo] {
	return this.todoRepository.FindOneWithItemsByID(id)
}

func (this *TodoServiceImpl) SaveList(todo *Todo) (*Todo, error) {
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()
	t, err := this.todoRepository.Save(todo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return t, nil
}

func (this *TodoServiceImpl) DeleteList(todo *Todo) error {
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()
	if err := this.todoRepository.DeleteById(todo.ID); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (this *TodoServiceImpl) ReorderLists(user *User, order []uint) error {
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()
	todos := this.todoRepository.FindAllByUserID(user.ID)
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
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()
	items := this.itemRepository.FindAllByTodoID(todo.ID)
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
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()
	item, ok := this.itemRepository.FindOneByTodoIDAndRank(todo.ID, oldRank).Get()
	if !ok {
		return NewError("itemRepositoryError", "Unable to find source item")
	}
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

func (this *TodoServiceImpl) GetTodoItem(todo *Todo, itemID uint) Optional[*Item] {
	return this.itemRepository.FindOneByTodoIDAndItemID(todo.ID, itemID)
}

func (this *TodoServiceImpl) GetItem(id uint) Optional[*Item] {
	return this.itemRepository.FindById(id)
}

func (this *TodoServiceImpl) GetItems(todo *Todo) []Item {
	return this.itemRepository.FindAllByTodoIDOrderByRank(todo.ID)
}

func (this *TodoServiceImpl) AddItem(todo *Todo) (*Item, error) {
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()
	var ids []uint
	var id uint = 1

	if todo == nil {
		err := NewError("TodoServiceError", "Todo is empty")
		err.LogError()
		return nil, err
	}

	this.itemRepository.Query(func(db *gorm.DB) *gorm.DB {
		return db.Where("todo_id = ?", todo.ID)
	}).Pluck("item_id", &ids)

	if len(ids) > 0 {
		id = slices.Max(ids) + 1
	}

	item_title := this.TranslateWithDefault("item.title", fmt.Sprintf("Item %d", id), strconv.FormatUint(uint64(id), 10))
	item := &Item{TodoID: todo.ID, ItemID: id, Text: item_title, Rank: uint(len(ids) + 1)}
	i, err := this.itemRepository.Save(item)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return i, nil
}

func (this *TodoServiceImpl) SaveItem(item *Item) (*Item, error) {
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()
	t, err := this.itemRepository.Save(item)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return t, nil
}

func (this *TodoServiceImpl) DeleteItem(item *Item) error {
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()
	if err := this.itemRepository.DeleteById(item.ID); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
