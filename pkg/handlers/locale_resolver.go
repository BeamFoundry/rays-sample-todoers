package handlers

import (
	. "github.com/BeamFoundry/rays/pkg/i18n"
	. "github.com/BeamFoundry/rays/pkg/i18n/locales"

	"github.com/gofiber/fiber/v3"
)

type LocaleResolver struct {
	localeContextHolder LocaleContextHolder `@:"Inject"`
}

func (this *LocaleResolver) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		if lang := c.Query("lang"); lang != "" {
			switch lang {
			case "fr":
				lang = "fr"
			default:
				lang = "en"
			}
			c.Cookie(&fiber.Cookie{
				Name:        "lang",
				Value:       lang,
				Domain:      c.Hostname(),
				SessionOnly: true,
			})
			return c.Redirect().Status(fiber.StatusFound).To(c.Path())
		}

		current_locale := "en"

		if lang := c.Cookies("lang"); lang != "" {
			switch lang {
			case "fr":
				lang = "fr"
			default:
				lang = "en"
			}
			current_locale = lang
		}

		fiber.Locals(c, "Lang", current_locale)
		switch current_locale {
		case "fr":
			this.localeContextHolder.Set(Locales.FR)
		default:
			this.localeContextHolder.Set(Locales.EN)
		}
		return c.Next()
	}
}
