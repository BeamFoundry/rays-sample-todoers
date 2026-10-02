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

func sub(f fs.FS, path string) fs.FS {
	subFS, _ := fs.Sub(f, path)
	return subFS
}

func Config() fs.FS {
	return config
}

func Views() fs.FS {
	return views
}

func Static() fs.FS {
	return sub(static, "static")
}
