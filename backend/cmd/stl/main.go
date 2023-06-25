package main

import (
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

func main() {
	app := a.NewApp(name, env, log)

	err := app.Run()
	if err != nil {
		log.Errorf("%s exit error: %s", name, err.Error())
		os.Exit(1)
	}
}
