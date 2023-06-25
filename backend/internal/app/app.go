package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/vanillazen/stl/backend/internal/domain/port"
	"github.com/vanillazen/stl/backend/internal/domain/service"
	"github.com/vanillazen/stl/backend/internal/infra/db"
	"github.com/vanillazen/stl/backend/internal/infra/db/sqlite"
	http2 "github.com/vanillazen/stl/backend/internal/infra/http"
	sqlite2 "github.com/vanillazen/stl/backend/internal/infra/sqlite"

	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/config"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
	"github.com/vanillazen/stl/backend/internal/sys/log"
)

type App struct {
	sync.Mutex
	sys.Core
	opts       []sys.Option
	supervisor sys.Supervisor
	http       *http2.Server
	db         db.DB
	repo       port.Repo
	svc        service.ListService
}

func NewApp(name, namespace string, log log.Logger) (app *App) {
	cfg := config.Load(namespace)

	opts := []sys.Option{
		sys.WithConfig(cfg),
		sys.WithLogger(log),
	}

	app = &App{
		Core: sys.NewCore(name, opts...),
		opts: opts,
	}

	return app
}

func (app *App) Run() (err error) {
	ctx := context.Background()

	err = app.Setup(ctx)
	if err != nil {
		return errors.Wrap(runError, err)
	}

	return app.Start(ctx)
}

func (app *App) Setup(ctx context.Context) error {
	// Databases
	dbase := sqlite.NewDB(app.opts...)

	// Repos
	repo, err := sqlite2.NewListRepo(dbase, app.opts...)
	if err != nil {
		return err
	}

	// Services
	app.svc = service.NewService(repo, app.opts...)

	// HTTP Server
	app.http = http2.NewServer(app.svc, app.opts...)

	return nil
}

func (app *App) Start(ctx context.Context) error {
	app.Log().Infof("%s starting...", app.Name())
	defer app.Log().Infof("%s stopped", app.Name())

	err := app.db.Start(ctx)
	if err != nil {
		msg := fmt.Sprintf("%s start error", app.db.Name())
		return errors.Wrap(msg, err)
	}

	err = app.svc.Start(ctx)
	if err != nil {
		return errors.Wrap("app start error", err)
	}

	app.supervisor.AddTasks(
		app.http.Start,
		//app.grpc.Start,
	)

	app.Log().Infof("%s started!", app.Name())

	return app.supervisor.Wait()
}

func (app *App) Stop(ctx context.Context) error {
	return nil
}

func (app *App) Shutdown(ctx context.Context) error {
	return nil
}
