package controllers

import (
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/controllers/handlers"
)

type GlobalPathFilter struct {
	FiberPathFilter
	Compress
	Rewrite
	Logger
	LocaleResolver
}
