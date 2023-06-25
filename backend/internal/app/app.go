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
	sqliterepo "github.com/vanillazen/stl/backend/internal/infra/repo/sqlite"

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
	repo       port.ListRepo
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
	app.EnableSupervisor()

	// Databases
	app.db = sqlite.NewDB(app.opts...)

	// Repos
	app.repo = sqliterepo.NewListRepo(app.db, app.opts...)

	// Services
	app.svc = service.NewService(app.repo, app.opts...)

	// HTTP Server
	app.http = http2.NewServer(app.svc, app.opts...)

	return nil
}

func (app *App) Start(ctx context.Context) error {
	app.Log().Infof("%s starting...", app.Name())
	defer app.Log().Infof("%s stopped", app.Name())

	app.supervisor.AddTasks(
		//app.db.Start,
		app.repo.Start,
		app.svc.Start,
		app.http.Start,
		//app.grpc.Start,
	)

	app.supervisor.AddShutdownTasks(
		app.http.Stop,
		//app.grpc.Start,
		app.svc.Stop,
		app.repo.Stop,
		//app.db.Stop,
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

func (app *App) EnableSupervisor() {
	name := fmt.Sprintf("%s-supervisor", app.Name())
	app.supervisor = sys.NewSupervisor(name, true, app.opts)
}
