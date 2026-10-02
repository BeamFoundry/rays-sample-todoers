package config

import (
	"errors"
	"net/http"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/handlers"

	"github.com/BeamFoundry/rays-sample-todoers/resources"
	"github.com/BeamFoundry/rays-fiber/pkg/template/html"

	"github.com/gofiber/fiber/v3"
)

type FiberAppConfigCustomizer struct {
	FiberPathFilter
	Compress
	Rewrite
	Logger
	LocaleResolver
	logger          Logging               `@:"Inject"`
	structValidator fiber.StructValidator `@:"Inject"`
	fnMaps          []FuncMapper          `@:"Inject"`
}

//nolint:errcheck // TODO : Add check of fnMaps.Load
func (this *FiberAppConfigCustomizer) FiberConfig(config *fiber.Config) {
	this.logger.Debug("Configuring a new fiber.Config")
	views := html.NewFileSystem(http.FS(resources.Views()), ".tmpl")
	fnMaps := TemplateFunctionMap{}.New()
	for _, fnMap := range this.fnMaps {
		fnMaps.Load(fnMap.FuncMap())
	}

	// TODO: Verify it can be called only once.
	views.AddFuncMap(fnMaps.FuncMap())

	errorHandler := func(c fiber.Ctx, err error) error {
		httpError := struct {
			Message string
			Code    int
		}{Code: fiber.StatusInternalServerError}

		// Retrieve the custom status code if it's a *fiber.Error
		if e, ok := errors.AsType[*fiber.Error](err); ok {
			httpError.Code = e.Code
			httpError.Message = e.Message
		}
		return c.Render("views/error", fiber.Map{"Error": httpError})
	}
	config.StructValidator = this.structValidator
	config.PassLocalsToViews = true
	config.ViewsLayout = "views/layouts/main"
	config.Views = views
	config.ErrorHandler = errorHandler
}
