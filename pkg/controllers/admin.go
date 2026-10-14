package controllers

import (
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/controllers/handlers"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/core"

	"github.com/gofiber/fiber/v3"
)

type AdminController struct {
	FiberController
	RequestMapping `Path:"/admin"`
	AppSecure
	userService UserService `@:"Inject"`
}

func (this *AdminController) AdminHome(
	_ struct {
		RequestMapping `Path:"/"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.Render("views/admin", fiber.Map{}, "views/layouts/admin")
}

func (this *AdminController) AdminListUsers(
	_ struct {
		RequestMapping `Path:"/users"`
	},
	ctx fiber.Ctx,
) error {
	page := fiber.Query(ctx, "page", 1)
	query := ctx.Query("q", "")
	count := this.userService.CountAll()
	users := this.userService.FindAllUsersByUsername(query, page, 10)

	m := fiber.Map{
		"CurrentPage": page,
		"SearchQuery": query,
		"Users":       users,
	}
	if count > int64(page*10) {
		m["NextPage"] = page + 1
	}
	return ctx.Render("views/fragments/admin-user-list", m)
}
