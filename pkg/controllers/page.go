package controllers

import (
	"errors"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-web-security/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/i18n"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"

	"github.com/go-playground/validator/v10"
)

type PageController struct {
	FiberController
	securityContextHolder SecurityContextHolder `@:"Inject"`

	userService            UserService            `@:"Inject"`
	messageSourceValidator MessageSourceValidator `@:"Inject"`
}

func (this *PageController) Home(
	_ struct {
		RequestMapping `Path:"/"`
	},
	c fiber.Ctx,
) error {
	auth := this.securityContextHolder.GetContext().GetAuthentication()
	if !auth.IsAuthenticated() {
		return c.Render("views/index", fiber.Map{})
	} else if user, ok := auth.GetPrincipal().(*User); !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	} else {
		values := fiber.Map{"User": user}
		sess := session.FromContext(c)
		if message, ok := sess.Get("message").(string); ok {
			sess.Delete("message")
			values["Message"] = message
		}
		return c.Render("views/app", values, "views/layouts/application")
	}
}

func (this *PageController) ViewProfile(
	_ struct {
		RequestMapping `Path:"/profile"`
	},
	c fiber.Ctx,
) error {
	auth := this.securityContextHolder.GetContext().GetAuthentication()
	if user, ok := auth.GetPrincipal().(*User); !ok {
		return c.SendStatus(fiber.StatusForbidden)
	} else {
		return c.Render("views/profile", fiber.Map{"User": user})
	}
}

func (this *PageController) EditProfile(
	_ struct {
		RequestMapping `Method:"POST" Path:"/profile"`
	},
	c fiber.Ctx,
) error {
	type Profile struct {
		Name     string `form:"name" validate:"required"`
		Password string `form:"password" validate:"required_with=Confirm"`
		Confirm  string `form:"confirm" validate:"eqfield=Password"`
	}
	form := Profile{}

	auth := this.securityContextHolder.GetContext().GetAuthentication()
	messages := []string{}
	if user, ok := auth.GetPrincipal().(*User); !ok {
		return c.SendStatus(fiber.StatusForbidden)
	} else if err := c.Bind().Form(&form); err != nil {
		if errs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			return c.Render("views/profile", fiber.Map{"User": user, "Form": form, "Error": err.Error()})
		} else {
			errs := this.messageSourceValidator.TranslateValidationErrors(errs)
			return c.Render("views/profile", fiber.Map{"User": user, "Form": form, "Errors": errs})
		}
	} else {
		if form.Name != user.Name {
			if u, err := this.userService.UpdateName(form.Name, user); err != nil {
				return c.Render("views/profile", fiber.Map{"User": user, "Error": err})
			} else {
				user = u
				messages = append(messages, "profile.informations.updated")
			}
		}
		if form.Password != "" {
			if u, err := this.userService.UpdatePassword(form.Password, user); err != nil {
				return c.Render("views/profile", fiber.Map{"User": user, "Error": err})
			} else {
				user = u
				messages = append(messages, "profile.password.update")
			}
		}
		return c.Render("views/profile", fiber.Map{"User": user, "Messages": messages})
	}
}
