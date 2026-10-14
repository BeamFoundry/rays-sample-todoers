package resources

import (
	"embed"
	"io/fs"
)

//go:embed config
var config embed.FS

//go:embed static
var static embed.FS

//go:embed views
var views embed.FS

func Config() fs.FS {
	return config
}

func Views() fs.FS {
	return views
}

func Static() fs.FS {
	//nolint:errcheck // embed
	subFS, _ := fs.Sub(static, "static")
	return subFS
}
