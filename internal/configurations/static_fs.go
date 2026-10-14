package configurations

import (
	"io/fs"

	"github.com/BeamFoundry/rays-sample-todoers/resources"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	"github.com/benbjohnson/hashfs"
)

type HashFS struct {
	Component
	hashFS *hashfs.FS
}

func (this HashFS) New() *HashFS {
	this.hashFS = hashfs.NewFS(resources.Static())
	return &this
}

func (this *HashFS) HashFS() *hashfs.FS {
	return this.hashFS
}

func (this *HashFS) StaticFS() fs.FS {
	return this.hashFS
}

var _ StaticFS = &HashFS{}

type HashFSFuncConfig struct {
	Configuration
	hashFS *HashFS `@:"Inject"`
}

func (this *HashFSFuncConfig) HashFSFunc() (FuncMapper, Error) {
	fnmap := TemplateFunctionMap{}.New()
	if err := fnmap.Add("hashfs", func(name string) string {
		return this.hashFS.HashFS().HashName(name)
	}); err != nil {
		return nil, err
	}

	return fnmap, nil
}
