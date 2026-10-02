package controllers

import (
	"html/template"
	"io/fs"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	"github.com/BeamFoundry/rays-sample-todoers/resources"
	"github.com/benbjohnson/hashfs"
)

type StaticController struct {
	Controller
	RequestMapping `Path:"/static"`
	hashFS         *hashfs.FS
	Etag
	Static
}

var _ StaticFS = &StaticController{}

func (this StaticController) New() *StaticController {
	this.hashFS = hashfs.NewFS(resources.Static())
	return &this
}

func (this *StaticController) StaticFS() fs.FS {
	return this.hashFS
}

//nolint:errcheck // TODO: Evaluate if error check is needed
func (this *StaticController) FuncMap() template.FuncMap {
	fnmap := TemplateFunctionMap{}.New()
	fnmap.Add("hashfs", func(name string) string {
		return this.hashFS.HashName(name)
	})

	return fnmap.FuncMap()
}
