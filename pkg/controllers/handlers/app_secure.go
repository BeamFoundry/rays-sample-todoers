package handlers

import (
	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"

	"github.com/gofiber/fiber/v3"
)

type AppSecure struct {
	securityContextHolder SecurityContextHolder `@:"Inject"`
}

func (this *AppSecure) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		auth := this.securityContextHolder.GetContext().GetAuthentication()
		if !auth.IsAuthenticated() {
			return c.Redirect().Status(302).To("/register")
		}
		if u, ok := auth.GetPrincipal().(*User); !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		} else {
			fiber.Locals(c, "User", u)
		}
		return c.Next()
	}
}

type AppSecureTodo struct {
	todoService TodoService `@:"Inject"`
}

func (this *AppSecureTodo) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		todo_id := fiber.Params[uint](c, "id", 0)
		if todo_id == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		user := fiber.Locals[*User](c, "User")
		if todo, ok := this.todoService.GetUserList(user, todo_id).Get(); !ok {
			return c.SendStatus(fiber.StatusBadRequest)
		} else if todo.UserID != user.ID {
			return c.SendStatus(fiber.StatusForbidden)
		} else {
			fiber.Locals(c, "Todo", todo)
			fiber.Locals(c, "TodoID", todo.TodoID)
		}
		return c.Next()
	}
}

type AppSecureItem struct {
	todoService TodoService `@:"Inject"`
}

func (this *AppSecureItem) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		item_id := fiber.Params[uint](c, "item", 0)
		if item_id == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		todo := fiber.Locals[*Todo](c, "Todo")
		if item, ok := this.todoService.GetTodoItem(todo, item_id).Get(); !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		} else if item.TodoID != todo.ID {
			return c.SendStatus(fiber.StatusBadRequest)
		} else {
			fiber.Locals(c, "Item", item)
		}
		return c.Next()
	}
}
