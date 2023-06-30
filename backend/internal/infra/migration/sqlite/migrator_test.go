package sqlite

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/vanillazen/stl/backend/internal/infra/db/sqlite"
	migrator "github.com/vanillazen/stl/backend/internal/infra/migration"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/config"
	l "github.com/vanillazen/stl/backend/internal/sys/log"
	"github.com/vanillazen/stl/backend/internal/sys/test"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

const (
	tmpDir     = "tmp"
	dbFileName = "stl-test.db"
	migDir     = "migrations"
)

var (
	logger l.Logger
	cfg    *config.Config
	opts   []sys.Option
	key    = config.Key
	testDB *sqlite.DB
)

type (
	TestCase struct {
		name     string
		cfg      *config.Config
		log      l.Logger
		opts     []sys.Option
		tmpPath  string
		fs       embed.FS
		db       *sqlite.DB
		migrator migrator.Migrator
		testFunc func(t *testing.T)
		expected test.Result
		result   test.Result
	}

	Result struct {
		value interface{}
		err   error
	}
)

func TestMain(m *testing.M) {
	setup()
	ev := m.Run()
	teardown()

	os.Exit(ev)
}

func TestMigrate(t *testing.T) {
	var tcs test.Cases

	tc0 := &TestCase{
		name:     "TestMigrateHappyPath",
		expected: Result{},
	}
	tc0.testFunc = tc0.TestMigrateHappyPath
	tcs.Add(tc0)

	//tc1 := &TestCase{
	//	name:     "TestMigrateCond0",
	//	expected: Result{},
	//}
	//tc1.testFunc = tc1.TestMigrateCond0
	//tcs.Add(tc1)

	// add more test cases...

	tests := tcs.All()

	for i := range tcs.All() {
		tc := tests[i]

		err := tc.Setup()
		if err != nil {
			t.Fatalf("%s setup error: %s", tc.Name(), err)
		}

		t.Run(tc.Name(), func(t *testing.T) {
			tc.TestFunc(t)
		})

		resErr := tc.Expected().Error()
		expErr := tc.Result().Error()

		resVal := tc.Expected().Value()
		expVal := tc.Result().Value()

		if resErr != expErr {
			t.Errorf("expected error '%s' but got: '%s'", expErr, resErr)

		} else if reflect.DeepEqual(resVal, expVal) {
			t.Errorf("expected value '%s' but got: '%s'", expVal, resVal)
		} else {
			t.Logf("%s: OK", tc.Name())
		}

		err = tc.Teardown()
		if err != nil {
			t.Errorf("%s teardown error: %s", tc.Name(), err)
		}
	}
}

func (tc *TestCase) TestMigrateHappyPath(t *testing.T) {
	tc.migrator = NewMigrator(tc.fs, tc.db)

	err := tc.migrator.Migrate()
	if err != nil {
		tc.result = Result{
			err: err,
		}
		return
	}

	// TODO: add assertions

	// ...

	tc.result = Result{
		value: true, // TODO: Add proper result value if required
		err:   nil,
	}
}

func (tc *TestCase) TestMigrateCond0(t *testing.T) {
	tc.migrator = NewMigrator(tc.fs, tc.db)

	err := tc.migrator.Migrate()
	if err != nil {
		tc.result = Result{
			err: err,
		}
		return
	}

	// TODO: add assertions
	// ...

	tc.result = Result{
		value: true, // TODO: Add proper result value if required
		err:   nil,
	}
}

func TestRollback(t *testing.T) {
	var tcs []*TestCase

	for i := range tcs {
		tc := *tcs[i]

		err := tc.Setup()
		if err != nil {
		}

		t.Run(tc.Name(), func(t *testing.T) {
			tc.TestFunc(t)
		})

		resErr := tc.Expected().Error()
		expErr := tc.result.Error()

		resVal := tc.Expected().Value()
		expVal := tc.Result().Value()

		if resErr != expErr {
			t.Errorf("expected error '%s' but got: '%s'", expErr, resErr)

		} else if reflect.DeepEqual(resVal, expVal) {
			t.Errorf("expected value '%s' but got: '%s'", expVal, resVal)
		} else {
			t.Logf("%s: OK", tc.Name())
		}

		err = tc.Teardown()
		if err != nil {
			t.Errorf("%s teardown error: %s", tc.Name(), err)
		}
	}
}

func (tc *TestCase) Name() string {
	return tc.name
}

func (tc *TestCase) Expected() test.Result {
	return tc.expected
}

func (tc *TestCase) Result() test.Result {
	return tc.result
}

func (tc *TestCase) TestFunc(t *testing.T) func(t *testing.T) {
	return tc.testFunc
}

func (tc *TestCase) Setup() error {
	ctx := context.Background()

	tc.cfg = cfg
	tc.log = logger
	tc.opts = opts
	tc.fs = embed.FS{}
	tc.db = sqlite.NewDB(tc.opts...)

	// Create the temporary directory
	tmpPath := filepath.Join(".", tmpDir)
	err := os.MkdirAll(tmpPath, 0755)
	if err != nil {
		msg := fmt.Errorf("failed to create tmp dir: %v", err)
		panic(msg)
	}

	// Create the migrations directory
	migrationsPath := filepath.Join(tmpPath, migDir)
	err = os.Mkdir(migrationsPath, 0755)
	if err != nil {
		err := fmt.Errorf("failed to create temp migrations dir: %v", err)
		return err
	}

	// Set config values to test temp directories
	opts := cfg.GetValues()
	opts[key.SQLiteFilePath] = filepath.Join(tmpPath, dbFileName)
	cfg.SetValues(opts)

	// Create DB
	// DB will be used for assertions for this reason we maintain a global reference.
	// In SQLite we don't want to have multiple open connections.
	testDB = sqlite.NewDB(tc.opts...)
	err = testDB.Start(ctx)
	if err != nil {
		err := fmt.Errorf("failed to create test db: %v", err)
		return err
	}

	tc.db = testDB

	return nil
}

func (tc *TestCase) Teardown() error {
	ctx := context.Background()
	err := tc.db.Stop(ctx)
	if err != nil {
		err = fmt.Errorf("failed to remove tmp dir: %v", err)
		logger.Error(err)
	}

	tmpPath := filepath.Join(".", tmpDir)

	err = os.RemoveAll(tmpPath)
	if err != nil {
		err := fmt.Errorf("failed to remove tmp dir: %v", err)
		logger.Error(err)
	}

	return nil
}

func setup() {
	cfg = &config.Config{}
	cfg.SetValues(optValues())

	// Test logger
	logger = l.NewTestLogger("error")

	// opts
	opts = []sys.Option{
		sys.WithConfig(cfg),
		sys.WithLogger(logger),
	}
}

// teardown removes the temporary directory and files created for the tests
func teardown() {
}

func optValues() map[string]string {
	return map[string]string{
		key.SQLiteUser:   "stl",
		key.SQLitePass:   "stl",
		key.SQLiteDB:     "stl-test",
		key.SQLiteHost:   "localhost",
		key.SQLitePort:   "",
		key.SQLiteSchema: "",
		key.SQLiteSSL:    "false",
	}
}

func (r Result) Value() interface{} {
	return r.value
}

func (r Result) Error() error {
	return r.err
}
