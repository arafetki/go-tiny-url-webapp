package assets

import (
	"embed"
)

//go:embed "migrations"
var MigrationFiles embed.FS

//go:embed "static"
var StaticFiles embed.FS

//go:embed "templates"
var Templates embed.FS
