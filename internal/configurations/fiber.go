package configurations

import (
	"errors"
	"net/http"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/web"

	"github.com/BeamFoundry/rays-fiber/pkg/template/html"
	"github.com/BeamFoundry/rays-sample-todoers/resources"

	"github.com/gofiber/fiber/v3"
)

type FiberAppConfigCustomizer struct {
	Component
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
		}{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}
		// Retrieve the custom status code if it's a *fiber.Error
		if e, ok := errors.AsType[*fiber.Error](err); ok {
			httpError.Code = e.Code
			httpError.Message = e.Message
		} else if e, ok := errors.AsType[Error](err); ok {
			e.LogErrorWithStackTrace()
		}
		return c.Render("views/error", fiber.Map{"Error": httpError})
	}
	config.StructValidator = this.structValidator
	config.PassLocalsToViews = true
	config.ViewsLayout = "views/layouts/main"
	config.Views = views
	config.ErrorHandler = errorHandler
}
