package app

import (
	"context"
	"embed"
	"fmt"
	"sync"

	"github.com/vanillazen/stl/backend/internal/domain/port"
	"github.com/vanillazen/stl/backend/internal/domain/service"
	"github.com/vanillazen/stl/backend/internal/infra/db"
	"github.com/vanillazen/stl/backend/internal/infra/db/sqlite"
	"github.com/vanillazen/stl/backend/internal/infra/fixture"
	sqlite2 "github.com/vanillazen/stl/backend/internal/infra/fixture/sqlite"
	http2 "github.com/vanillazen/stl/backend/internal/infra/http"
	migrator "github.com/vanillazen/stl/backend/internal/infra/migration"
	mig "github.com/vanillazen/stl/backend/internal/infra/migration/sqlite"
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
	fs         embed.FS
	supervisor sys.Supervisor
	http       *http2.Server
	db         db.DB
	repo       port.ListRepo
	migrator   migrator.Migrator
	fixture    fixture.Fixture
	svc        service.ListService
	apiDoc     string
}

func NewApp(name, namespace string, apiDoc string, fs embed.FS, log log.Logger) (app *App) {
	cfg := config.Load(namespace)

	opts := []sys.Option{
		sys.WithConfig(cfg),
		sys.WithLogger(log),
	}

	app = &App{
		Core:   sys.NewCore(name, opts...),
		opts:   opts,
		fs:     fs,
		apiDoc: apiDoc,
	}

	return app
}

func (app *App) Run() (err error) {
	ctx := context.Background()

	err = app.Setup(ctx)
	if err != nil {
		return errors.Wrap(err, runError)
	}

	return app.Start(ctx)
}

func (app *App) Setup(ctx context.Context) error {
	app.EnableSupervisor()

	// Databases
	app.db = sqlite.NewDB(app.opts...)

	// Migration
	app.migrator = mig.NewMigrator(app.fs, app.db, app.opts...)

	// Pre-population
	app.fixture = sqlite2.NewFixture(app.db, app.opts...)

	// Repos
	app.repo = sqliterepo.NewListRepo(app.db, app.opts...)

	// Services
	app.svc = service.NewService(app.repo, app.opts...)

	// HTTP Server
	app.http = http2.NewServer(app.svc, app.apiDoc, app.opts...)

	err := app.http.Setup(ctx)
	if err != nil {
		err = errors.Wrapf(err, "%s setup error", app.Name())
		return err
	}

	err = app.db.Start(ctx)
	if err != nil {
		err = errors.Wrapf(err, "%s setup error", app.Name())
		return err
	}

	return nil
}

func (app *App) Start(ctx context.Context) (err error) {
	app.Log().Infof("%s starting...", app.Name())
	defer app.Log().Infof("%s stopped", app.Name())

	//// Non-blocking sequential start
	err = app.repo.Start(ctx)
	if err != nil {
		app.Log().Errorf("%s start error: %s", err)
		return err
	}

	err = app.migrator.Start(ctx)
	if err != nil {
		app.Log().Errorf("%s start error: %s", app.Name(), err)
		return err
	}

	err = app.fixture.Start(ctx)
	if err != nil {
		err = errors.Wrapf(err, "%s setup error", app.Name())
		return err
	}

	err = app.svc.Start(ctx)
	if err != nil {
		app.Log().Errorf("%s start error: %s", app.Name(), err)
		return err
	}

	// Blocking non-sequential start
	app.supervisor.AddTasks(
		app.http.Start,
		//app.grpc.Start,
	)

	app.supervisor.AddShutdownTasks(
		app.http.Stop,
		//app.grpc.Start,
	)

	app.Log().Infof("%s started", app.Name())

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
