package sqlite

import (
	"embed"
	"os"
	"reflect"
	"testing"

	"github.com/vanillazen/stl/backend/internal/infra/db"
	"github.com/vanillazen/stl/backend/internal/infra/db/sqlite"
	migrator "github.com/vanillazen/stl/backend/internal/infra/migration"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/config"
	l "github.com/vanillazen/stl/backend/internal/sys/log"
	"github.com/vanillazen/stl/backend/internal/sys/test"
)

type (
	TestCase struct {
		name     string
		cfg      *config.Config
		log      l.Logger
		opts     []sys.Option
		fs       embed.FS
		db       db.DB
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

var (
	logger l.Logger
	cfg    *config.Config
	opts   []sys.Option
	FS     embed.FS
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

	tc1 := &TestCase{
		name:     "TestMigrateCond0",
		expected: Result{},
	}
	tc1.testFunc = tc1.TestMigrateCond0

	// add more test cases...
	tcs.Add(tc0, tc1)

	tests := tcs.All()

	for i := range tcs.All() {
		tc := tests[i]

		err := tc.Setup()
		if err != nil {
			t.Errorf("%s setup error: %s", tc.Name(), err)
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
	tc.cfg = cfg
	tc.log = logger
	tc.opts = opts
	tc.fs = FS
	tc.db = sqlite.NewDB(tc.opts...)

	return nil
}

func (tc *TestCase) Teardown() error {
	// TODO: Cleanup resources
	return nil
}

func (r Result) Value() interface{} {
	return r.value
}

func (r Result) Error() error {
	return r.err
}

func setup() {
	cfg = &config.Config{}
	cfg.SetValues(map[string]string{}) // TODO: Set required config values

	// Test logger
	logger = l.NewTestLogger("debug")

	// opts
	opts = []sys.Option{
		sys.WithConfig(cfg),
		sys.WithLogger(logger),
	}

	// Embed filesystem
	FS = embed.FS{} // TODO: mock fs
}

func teardown() {
	// TODO: General test suite teardown
}
