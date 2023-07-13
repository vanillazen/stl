package sqlite

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/vanillazen/stl/backend/internal/infra/db"
	"github.com/vanillazen/stl/backend/internal/infra/db/sqlite"
	migrator "github.com/vanillazen/stl/backend/internal/infra/migration"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/config"
	l "github.com/vanillazen/stl/backend/internal/sys/log"
	"github.com/vanillazen/stl/backend/internal/sys/test"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tmpDir     = "tmp"
	dbFileName = "stl-test.db"
)

var (
	logger l.Logger
	cfg    *config.Config
	opts   []sys.Option

	assetsPath  = filepath.Join(tmpDir, "seeding")
	assetsPath1 = filepath.Join(tmpDir, "seeding-1")

	//go:embed all:tmp/seeding/*.sql
	fs embed.FS

	//go:embed all:tmp/seeding-1/*.sql
	fs1 embed.FS

	key      = config.Key
	testDB   *sqlite.DB
	testsDir = []string{}
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

func NewTestCase(name string,
	cfg *config.Config,
	logger l.Logger,
	opts []sys.Option,
	tnpPath string) *TestCase {

	return &TestCase{
		name:     name,
		cfg:      cfg,
		log:      logger,
		opts:     opts,
		tmpPath:  tnpPath,
		fs:       embed.FS{},
		db:       nil,
		migrator: nil,
		testFunc: nil,
		expected: Result{
			value: nil,
			err:   nil,
		},
		result: Result{
			value: nil,
			err:   nil,
		},
	}
}

func TestMain(m *testing.M) {
	setup()

	ev := m.Run()

	teardown()

	os.Exit(ev)
}

func TestMigrator(t *testing.T) {
	var tcs test.Cases

	tc0 := NewTestCase("TestMigrateBase", cfg, logger, opts, tmpDir)
	tc0.testFunc = tc0.TestMigrateBase
	tc0.expected = Result{
		err: nil,
	}
	tcs.Add(tc0)

	tc1 := NewTestCase("TestMigrateAndAgain", cfg, logger, opts, tmpDir)
	tc1.testFunc = tc0.TestMigrateAndAgain
	tc1.expected = Result{
		err: nil,
	}
	tcs.Add(tc1)

	tc2 := NewTestCase("TestRollback1", cfg, logger, opts, tmpDir)
	tc2.testFunc = tc2.TestRollback1
	tc2.expected = Result{
		err: nil,
	}
	tcs.Add(tc2)

	tc3 := NewTestCase("TestRollback2", cfg, logger, opts, tmpDir)
	tc3.testFunc = tc3.TestRollback2
	tc3.expected = Result{
		err: nil,
	}
	tcs.Add(tc3)

	tc4 := NewTestCase("TestRollbackAll", cfg, logger, opts, tmpDir)
	tc4.testFunc = tc4.TestRollbackAll
	tc4.expected = Result{
		err: nil,
	}
	tcs.Add(tc4)

	// add more test cases...

	tests := tcs.All()

	for i := range tcs.All() {
		tc := tests[i]

		err := tc.Setup()
		if err != nil {
			t.Fatalf("%s setup error: %s", tc.Name(), err)
		}

		t.Run(tc.Name(), tc.TestFunc())

		resErr := tc.Expected().Error()
		expErr := tc.Result().Error()

		resVal := tc.Expected().Value()
		expVal := tc.Result().Value()

		if resErr != expErr {
			t.Errorf("expected error '%s' but got: '%s'", expErr, resErr)

		} else if !reflect.DeepEqual(resVal, expVal) {
			t.Errorf("expected value '%v' but got: '%v'", expVal, resVal)
		} else {
			t.Logf("%s: OK", tc.Name())
		}

		err = tc.Teardown()
		if err != nil {
			t.Errorf("%s teardown error: %s", tc.Name(), err)
		}
	}
}

func (tc *TestCase) TestMigrateBase(t *testing.T) {
	ctx := context.Background()
	tc.migrator = NewMigrator(tc.fs, tc.db, tc.opts...)
	tc.migrator.SetAssetsPath(assetsPath)

	err := tc.migrator.Start(ctx)
	if err != nil {
		tc.result = Result{
			err: err,
		}
		return
	}

	expected := []string{"users", "lists", "tasks"}
	found, ok := tablesExists(tc.db, expected...)
	if !ok {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	migRecords, err := migRecords(tc.db)
	if err != nil {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	expectedIdxSum := 6
	idxSum := 0
	names := []string{}
	validCreatedAt := true

	for _, mg := range migRecords {
		names = append(names, mg.Name)
		idxSum = idxSum + mg.Index
		validCreatedAt = validCreatedAt && !isAfterNowMinus(1, mg.CreatedAt)
	}

	if !reflect.DeepEqual(names, expected) {
		t.Errorf("expected migration record '%v' but got: '%v'", expected, names)
	}

	if idxSum != expectedIdxSum {
		t.Errorf("wrong migration record indexes")
	}

	if !validCreatedAt {
		t.Errorf("wrong migration record created at timestamp")
	}

	tc.result = Result{
		err: nil,
	}
}

func (tc *TestCase) TestMigrateAndAgain(t *testing.T) {
	ctx := context.Background()
	tc.migrator = NewMigrator(tc.fs, tc.db, tc.opts...)
	tc.migrator.SetAssetsPath(assetsPath)

	err := tc.migrator.Start(ctx)
	if err != nil {
		t.Error("pre-setup failed")
		tc.result = Result{
			err: err,
		}
		return
	}

	tc.fs = fs1
	tc.migrator = NewMigrator(tc.fs, tc.db, tc.opts...)
	tc.migrator.SetAssetsPath(assetsPath1)

	err = tc.migrator.Start(ctx)
	if err != nil {
		tc.result = Result{
			err: err,
		}
		return
	}

	expected := []string{"users", "lists", "tasks", "tags"}
	found, ok := tablesExists(tc.db, expected...)
	if !ok {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	migRecords, err := migRecords(tc.db)
	if err != nil {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	expectedIdxSum := 10
	idxSum := 0
	names := []string{}
	validCreatedAt := true

	for _, mg := range migRecords {
		names = append(names, mg.Name)
		idxSum = idxSum + mg.Index
		validCreatedAt = validCreatedAt && !isAfterNowMinus(1, mg.CreatedAt)
	}

	if !reflect.DeepEqual(names, expected) {
		t.Errorf("expected migration record '%v' but got: '%v'", expected, names)
	}

	if idxSum != expectedIdxSum {
		t.Errorf("wrong migration record indexes")
	}

	if !validCreatedAt {
		t.Errorf("wrong migration record created at timestamp")
	}

	tc.result = Result{
		err: nil,
	}
}

func (tc *TestCase) TestRollback1(t *testing.T) {
	ctx := context.Background()
	tc.migrator = NewMigrator(tc.fs, tc.db, tc.opts...)
	tc.migrator.SetAssetsPath(assetsPath)

	err := tc.migrator.Start(ctx)
	if err != nil {
		t.Error("pre-setup failed")

		tc.result = Result{
			err: err,
		}
		return
	}

	err = tc.migrator.Rollback(1)
	if err != nil {
		t.Errorf("error executing rollback: '%s'", err.Error())
	}

	expected := []string{"users", "lists"}
	found, ok := tablesExists(tc.db, expected...)
	if !ok {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	migRecords, err := migRecords(tc.db)
	if err != nil {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	expectedIdxSum := 3
	idxSum := 0
	names := []string{}
	validCreatedAt := true

	for _, mg := range migRecords {
		names = append(names, mg.Name)
		idxSum = idxSum + mg.Index
		validCreatedAt = validCreatedAt && !isAfterNowMinus(1, mg.CreatedAt)
	}

	if !reflect.DeepEqual(names, expected) {
		t.Errorf("expected migration record '%v' but got: '%v'", expected, names)
	}

	if idxSum != expectedIdxSum {
		t.Errorf("wrong migration record indexes")
	}

	if !validCreatedAt {
		t.Errorf("wrong migration record created at timestamp")
	}

	tc.result = Result{
		err: nil,
	}
}

func (tc *TestCase) TestRollback2(t *testing.T) {
	ctx := context.Background()
	tc.migrator = NewMigrator(tc.fs, tc.db, tc.opts...)
	tc.migrator.SetAssetsPath(assetsPath)

	err := tc.migrator.Start(ctx)
	if err != nil {
		t.Error("pre-setup failed")

		tc.result = Result{
			err: err,
		}
		return
	}

	err = tc.migrator.Rollback(2)
	if err != nil {
		t.Errorf("error executing rollback: '%s'", err.Error())
	}

	expected := []string{"users"}
	found, ok := tablesExists(tc.db, expected...)
	if !ok {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	migRecords, err := migRecords(tc.db)
	if err != nil {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	expectedIdxSum := 1
	idxSum := 0
	names := []string{}
	validCreatedAt := true

	for _, mg := range migRecords {
		names = append(names, mg.Name)
		idxSum = idxSum + mg.Index
		validCreatedAt = validCreatedAt && !isAfterNowMinus(1, mg.CreatedAt)
	}

	if !reflect.DeepEqual(names, expected) {
		t.Errorf("expected migration record '%v' but got: '%v'", expected, names)
	}

	if idxSum != expectedIdxSum {
		t.Errorf("wrong migration record indexes")
	}

	if !validCreatedAt {
		t.Errorf("wrong migration record created at timestamp")
	}

	tc.result = Result{
		err: nil,
	}
}

func (tc *TestCase) TestRollbackAll(t *testing.T) {
	ctx := context.Background()
	tc.migrator = NewMigrator(tc.fs, tc.db, tc.opts...)
	tc.migrator.SetAssetsPath(assetsPath)

	err := tc.migrator.Start(ctx)
	if err != nil {
		t.Error("pre-setup failed")

		tc.result = Result{
			err: err,
		}
		return
	}

	err = tc.migrator.RollbackAll()
	if err != nil {
		t.Errorf("error executing rollback: '%s'", err.Error())
	}

	expected := []string{}
	found, ok := tablesExists(tc.db, expected...)
	if !ok {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	migRecords, err := migRecords(tc.db)
	if err != nil {
		t.Errorf("expected tables '%v' but got: '%v'", expected, found)
	}

	expectedIdxSum := 0
	idxSum := 0
	names := []string{}
	validCreatedAt := true

	for _, mg := range migRecords {
		names = append(names, mg.Name)
		idxSum = idxSum + mg.Index
		validCreatedAt = validCreatedAt && !isAfterNowMinus(1, mg.CreatedAt)
	}

	if !reflect.DeepEqual(names, expected) {
		t.Errorf("expected migration record '%v' but got: '%v'", expected, names)
	}

	if idxSum != expectedIdxSum {
		t.Errorf("wrong migration record indexes")
	}

	if !validCreatedAt {
		t.Errorf("wrong migration record created at timestamp")
	}

	tc.result = Result{
		err: nil,
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

func (tc *TestCase) TestFunc() func(t *testing.T) {
	return tc.testFunc
}

func (tc *TestCase) Setup() error {
	ctx := context.Background()

	tc.cfg = cfg
	tc.log = logger
	tc.opts = opts
	tc.fs = fs
	tc.db = sqlite.NewDB(tc.opts...)

	// Create the temporary directory
	testDir, err := os.MkdirTemp(tmpDir, "test")
	if err != nil {
		msg := fmt.Errorf("failed to create tmp dir: %v", err)
		panic(msg)
	}

	testsDir = append(testsDir, testDir)

	// Set config values to test temp directories
	cfgValues := cfg.GetValues()
	cfgValues[key.SQLiteFilePath] = filepath.Join(testDir, dbFileName)
	cfg.SetValues(cfgValues)

	tc.opts = []sys.Option{
		sys.WithConfig(cfg),
		sys.WithLogger(logger),
	}

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
	for _, td := range testsDir {
		err := os.RemoveAll(td)
		if err != nil {
			logger.Error(err)
		}
	}
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

// Helpers

func migRecords(dbase db.DB) (migRecords []MigrationRecord, err error) {
	rows, err := dbase.DB().Query("SELECT id, idx, name, created_at FROM seeding")
	if err != nil {
		return migRecords, err
	}
	defer rows.Close()

	for rows.Next() {
		var record MigrationRecord

		// Scan the columns into the struct fields
		err := rows.Scan(&record.ID, &record.Index, &record.Name, &record.CreatedAt)
		if err != nil {
			return migRecords, err
		}

		migRecords = append(migRecords, record)
	}

	if err = rows.Err(); err != nil {
		return migRecords, err
	}

	return migRecords, nil
}

func tablesExists(dbase db.DB, tableNames ...string) (found []string, ok bool) {
	for _, t := range tableNames {
		var name string

		query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", t)

		err := dbase.DB().QueryRow(query).Scan(&name)
		if err != nil {
			continue
		}

		found = append(found, t)
	}

	return found, len(tableNames) == len(found)
}

func isAfterNowMinus(xMinutes int, dateStr string) (ok bool) {
	dateFmt := "2006-01-02 15:04:05"

	date, err := time.Parse(dateFmt, dateStr)
	if err != nil {
		return false
	}

	nowMinusXMinutes := time.Now().Add(time.Duration(-xMinutes) * time.Minute)

	if !date.After(nowMinusXMinutes) {
		return false
	}

	return true
}
