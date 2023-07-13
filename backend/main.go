package main

import (
	"embed"
	"os"

	a "github.com/vanillazen/stl/backend/internal/app"
	l "github.com/vanillazen/stl/backend/internal/sys/log"
)

const (
	name     = "stl"
	env      = "stl"
	logLevel = "info"
)

var (
	log l.Logger = l.NewLogger(logLevel)
)

var (
	//go:embed all:assets/migrations/sqlite/*.sql
	migFs embed.FS

	//go:embed all:assets/seeding/sqlite/*.sql
	seedFs embed.FS

	//go:embed api/openapi/openapi.html
	openapiDoc string
)

func main() {
	app := a.NewApp(name, env, log)
	app.SetMigratorFs(migFs)
	app.SetSeederFs(seedFs)
	app.SetAPIDoc(openapiDoc)

	err := app.Run()
	if err != nil {
		os.Exit(1)
	}
}
