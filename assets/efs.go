package assets

import (
	"embed"
)

//go:embed "migrations"
var Migrations embed.FS

//go:embed "static"
var Static embed.FS

//go:embed "templates"
var Templates embed.FS
