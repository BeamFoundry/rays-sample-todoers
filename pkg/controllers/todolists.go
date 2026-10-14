package controllers

import (
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/controllers/handlers"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"

	"github.com/gofiber/fiber/v3"
)

type TodolistsController struct {
	FiberController
	RequestMapping `Path:"/todolists"`
	AppSecure
	todoService TodoService `@:"Inject"`
}

func (this *TodolistsController) ListAll(
	_ struct {
		RequestMapping `Method:"GET"`
	},
	c fiber.Ctx,
) error {
	user := fiber.Locals[*User](c, "User")

	return c.Render("views/fragments/todolist", fiber.Map{
		"Todos": this.todoService.GetUserTodoLists(user),
	}, "views/fragments/todoitem")
}

func (this *TodolistsController) NewList(
	_ struct {
		RequestMapping `Method:"POST"`
	},
	c fiber.Ctx,
) error {
	user := fiber.Locals[*User](c, "User")
	todo, err := this.todoService.CreateTodolist(user)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.Render("views/fragments/todolist", fiber.Map{
		"Todos": []*Todo{todo},
	}, "views/fragments/todoitem")
}

func (this *TodolistsController) ReorderLists(
	_ struct {
		RequestMapping `Method:"PUT"`
	},
	c fiber.Ctx,
) error {
	user := fiber.Locals[*User](c, "User")

	type ReorderList struct {
		Todo []uint `form:"todo" validate:"required"`
	}

	form := ReorderList{}
	if err := c.Bind().Form(&form); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	} else if err := this.todoService.ReorderLists(user, form.Todo); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Render("views/fragments/todolist", fiber.Map{
		"Todos": this.todoService.GetUserTodoLists(user),
	}, "views/fragments/todoitem")
}

func (this *TodolistsController) ListsItemsActions(
	_ struct {
		RequestMapping `Method:"PUT" Path:"/:id/actions"`
		AppSecureTodo
	},
	c fiber.Ctx,
) error {
	todo := fiber.Locals[*Todo](c, "Todo")

	type ReorderListItem struct {
		Action string `form:"action" validate:"required"`
		Data   []uint `form:"data" validate:"required"`
	}

	form := ReorderListItem{}
	if err := c.Bind().Form(&form); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	switch form.Action {
	case "reorder":
		if err := this.todoService.ReorderListItems(todo, form.Data); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	case "move":
		if len(form.Data) != 3 {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := this.todoService.MoveListItem(todo, form.Data[0], form.Data[1], form.Data[2]); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (this *TodolistsController) EditList(
	_ struct {
		RequestMapping `Path:"/:id/edit"`
		AppSecureTodo
	},
	c fiber.Ctx,
) error {
	todo := fiber.Locals[*Todo](c, "Todo")
	return c.Render("views/fragments/edit-todolist-title", fiber.Map{
		"Todo": todo,
	})
}

func (this *TodolistsController) SaveList(
	_ struct {
		RequestMapping `Method:"PUT" Path:"/:id"`
		AppSecureTodo
	},
	c fiber.Ctx,
) error {
	todo := fiber.Locals[*Todo](c, "Todo")
	form := struct {
		Title string `form:"title" validate:"required"`
	}{}

	// TODO : Check again this works (list was named todo, shadowing todo)
	list, ok := this.todoService.GetListWithItems(todo.ID).Get()
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err := c.Bind().Form(&form); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	list.Title = form.Title
	t, err := this.todoService.SaveList(list)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Render("views/fragments/todolist", fiber.Map{
		"Todos": []*Todo{t},
	}, "views/fragments/todoitem")
}

func (this *TodolistsController) DeleteList(
	_ struct {
		RequestMapping `Method:"DELETE" Path:"/:id"`
		AppSecureTodo
	},
	c fiber.Ctx,
) error {
	todo := fiber.Locals[*Todo](c, "Todo")

	if err := this.todoService.DeleteList(todo); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString("")
}

func (this *TodolistsController) NewItem(
	_ struct {
		RequestMapping `Method:"POST" Path:"/:id/items"`
		AppSecureTodo
	},
	c fiber.Ctx,
) error {
	todo := fiber.Locals[*Todo](c, "Todo")
	item, err := this.todoService.AddItem(todo)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Render("views/fragments/todoitem", fiber.Map{
		"Items": []*Item{item},
	})
}

func (this *TodolistsController) EditItem(
	_ struct {
		RequestMapping `Path:"/:id/items/:item/edit"`
		AppSecureTodo
		AppSecureItem
	},
	c fiber.Ctx,
) error {
	item := fiber.Locals[*Item](c, "Item")
	return c.Render("views/fragments/edit-todoitem-text", fiber.Map{
		"Item": item,
	})
}

func (this *TodolistsController) SaveItem(
	_ struct {
		RequestMapping `Method:"PUT" Path:"/:id/items/:item"`
		AppSecureTodo
		AppSecureItem
	},
	c fiber.Ctx,
) error {
	form := struct {
		Text string `form:"text" validate:"required"`
	}{}

	item := fiber.Locals[*Item](c, "Item")
	if err := c.Bind().Form(&form); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	item.Text = form.Text
	i, err := this.todoService.SaveItem(item)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Render("views/fragments/todoitem", fiber.Map{
		"Items": []*Item{i},
	})
}

func (this *TodolistsController) ToggleItem(
	_ struct {
		RequestMapping `Method:"PUT" Path:"/:id/items/:item/done"`
		AppSecureTodo
		AppSecureItem
	},
	c fiber.Ctx,
) error {
	item := fiber.Locals[*Item](c, "Item")
	item.Done = !item.Done
	i, err := this.todoService.SaveItem(item)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Render("views/fragments/todoitem", fiber.Map{
		"Items": []*Item{i},
	})
}

func (this *TodolistsController) DeleteItem(
	_ struct {
		RequestMapping `Method:"DELETE" Path:"/:id/items/:item"`
		AppSecureTodo
		AppSecureItem
	},
	c fiber.Ctx,
) error {
	item := fiber.Locals[*Item](c, "Item")
	if err := this.todoService.DeleteItem(item); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString("")
}
