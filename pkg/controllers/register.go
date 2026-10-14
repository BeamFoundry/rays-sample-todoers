package controllers

import (
	"errors"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-web-security/pkg/auth"
	. "github.com/BeamFoundry/rays-web-security/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/i18n"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/core"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"

	"github.com/go-playground/validator/v10"
)

type RegisterController struct {
	FiberController
	RequestMapping         `Path:"/register"`
	securityContextHolder  SecurityContextHolder  `@:"Inject"`
	userService            UserService            `@:"Inject"`
	messageSourceValidator MessageSourceValidator `@:"Inject"`
}

func (this *RegisterController) Get(
	_ struct {
		RequestMapping `Method:"GET"`
	},
	c fiber.Ctx,
) error {
	auth := this.securityContextHolder.GetContext().GetAuthentication()
	if !auth.IsAuthenticated() {
		return c.Render("views/register", fiber.Map{})
	}
	return c.Redirect().Status(302).To("/")
}

func (this *RegisterController) Post(
	_ struct {
		RequestMapping `Method:"POST"`
	},
	c fiber.Ctx,
) error {
	type Register struct {
		Fullname string `form:"fullname" validate:"required"`
		Username string `form:"username" validate:"required"`
		Password string `form:"password" validate:"required"`
	}

	form := Register{}

	if err := c.Bind().Form(&form); err != nil {
		if errs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			tErrs := this.messageSourceValidator.TranslateValidationErrors(errs, "fields.register")
			return c.Render("views/register", fiber.Map{"Form": form, "Errors": tErrs})
		}
		return c.Render("views/register", fiber.Map{"Form": form, "Error": err.Error()})
	}

	u, err := this.userService.AddUser(form.Fullname, form.Username, form.Password)
	if err != nil {
		return c.Render("views/register", fiber.Map{"Error": err})
	}
	sess := session.FromContext(c)
	token := UsernamePasswordAuthenticationToken{}.NewAuthenticated(u, "", u.GetAuthorities())
	if err := sess.Reset(); err == nil {
		sess.Set("authentication", token)
		sess.Set("message", "register.successful")
		return c.Redirect().Status(302).To("/")
	}
	return nil
}
