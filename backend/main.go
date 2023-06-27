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

//go:embed all:assets/migrations/sqlite/*.sql
var fs embed.FS

func main() {
	app := a.NewApp(name, env, fs, log)

	err := app.Run()
	if err != nil {
		log.Errorf("%s exit error: %s", app.Name(), err.Error())
		os.Exit(1)
	}
}
