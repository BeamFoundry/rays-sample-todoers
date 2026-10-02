package config

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/web"
	"github.com/Masterminds/sprig/v3"
)

type SprigConfig struct {
	Configuration
}

func (this *SprigConfig) Sprig() (FuncMapper, Error) {
	fnmap := TemplateFunctionMap{}.New()
	if err := fnmap.Load(sprig.FuncMap()); err != nil {
		return nil, err
	}
	return fnmap, nil
}
