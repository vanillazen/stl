package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
	"github.com/vanillazen/stl/backend/internal/sys/uuid"
)

type (
	// Fx type alias
	Fx = func() error

	// Migrator struct.
	Migrator struct {
		sys.Core
		db     DB
		schema string
		dbPath string
		migs   []*Migration
	}

	// Exec interface.
	Exec interface {
		Config(up Fx, down Fx)
		GetName() (name string)
		GetUp() (up Fx)
		GetDown() (down Fx)
		SetTx(tx *sql.Tx)
		GetTx() (tx *sql.Tx)
	}

	// Migration struct.
	Migration struct {
		Order    int
		Executor Exec
	}

	migRecord struct {
		ID        uuid.UUID      `dbPath:"id" json:"id"`
		Name      sql.NullString `dbPath:"name" json:"name"`
		UpFx      sql.NullString `dbPath:"up_fx" json:"upFx"`
		DownFx    sql.NullString `dbPath:"down_fx" json:"downFx"`
		IsApplied sql.NullBool   `dbPath:"is_applied" json:"isApplied"`
		CreatedAt NullTime       `dbPath:"created_at" json:"createdAt"`
	}
)

const (
	sqlMigrationsTable = "migrations"

	sqlCreateDbSt = `
		CREATE DATABASE %s;`

	sqlDropDbSt = `
		DROP DATABASE %s;`

	sqlCreateMigrationsSt = `CREATE TABLE %s.%s (
		id UUID PRIMARY KEY,
		name VARCHAR(64),
		up_fx VARCHAR(64),
		down_fx VARCHAR(64),
 		is_applied BOOLEAN,
		created_at TIMESTAMP
	);`

	sqlDropMigrationsSt = `DROP TABLE %s.%s;`

	sqlSelMigrationSt = `SELECT is_applied FROM %s.%s WHERE name = '%s' and is_applied = true`

	sqlRecMigrationSt = `INSERT INTO %s.%s (id, name, up_fx, down_fx, is_applied, created_at)
		VALUES (:id, :name, :up_fx, :down_fx, :is_applied, :created_at);`

	sqlDelMigrationSt = `DELETE FROM %s.%s WHERE name = '%s' and is_applied = true`
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func NewMigrator(opts ...sys.Option) (mig *Migrator) {
	return &Migrator{
		Core: sys.NewCore("migrator", opts...),
	}
}

func (m *Migrator) DB() *sql.DB {
	return m.db.db
}

func (m *Migrator) Start(ctx context.Context) error {
	m.Log().Infof("%s started", m.db.Name())
	return m.Connect()
}

func (m *Migrator) Connect() error {
	sqlDB, err := sql.Open("sqlite3", m.db.dbPath())
	if err != nil {
		msg := fmt.Sprintf("%s connection error", m.db.Name())
		return errors.Wrap(err, msg)
	}

	err = sqlDB.Ping()
	if err != nil {
		msg := fmt.Sprintf("%s ping connection error", m.db.Name())
		return errors.Wrap(err, msg)
	}

	m.db.db = sqlDB
	m.Log().Infof("%s database connected", m.db.Name())
	return nil
}

// GetTx returns a new transaction from migrator connection
func (m *Migrator) GetTx() (tx *sql.Tx, err error) {
	tx, err = m.db.DB().Begin()
	if err != nil {
		return tx, err
	}

	return tx, nil
}

// PreSetup creates database
// and migrations table if needed.
func (m *Migrator) PreSetup() (err error) {
	if !m.dbExists() {
		_, err = m.CreateDb()
		if err != nil {
			return err
		}
	}

	if !m.migTableExists() {
		_, err = m.createMigrationsTable()
		if err != nil {
			return err
		}
	}

	return nil
}

// dbExists returns true if migrator referenced database has been already created.
func (m *Migrator) dbExists() bool {
	st := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='database' AND name='%s';", m.dbPath)

	rows, err := m.DB().Query(st)
	if err != nil {
		m.Log().Infof("Error checking database: %w", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var dbName string
		err = rows.Scan(&dbName)
		if err != nil {
			m.Log().Errorf("Cannot read query result: %w", err)
			return false
		}
		return true
	}

	return false
}

// migTableExists returns true if migrations table exists.
func (m *Migrator) migTableExists() bool {
	st := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", m.dbPath)

	rows, err := m.DB().Query(st)
	if err != nil {
		m.Log().Errorf("Error checking database: %s", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			m.Log().Errorf("Cannot read query result: %s\n", err)
			return false
		}

		return true
	}

	return false
}

// CreateDb migration.
func (m *Migrator) CreateDb() (dbPath string, err error) {
	_, err = m.CloseAppConns()
	if err != nil {
		return dbPath, errors.Wrap(err, "create db error")
	}

	st := fmt.Sprintf(sqlCreateMigrationsSt, m.dbPath)

	_, err = m.DB().Exec(st)
	if err != nil {
		return m.dbPath, err
	}

	return m.dbPath, nil
}

// DropDb migration.
func (m *Migrator) DropDb() (dbPath string, err error) {
	dbPath, err = m.CloseAppConns()
	if err != nil {
		return dbPath, errors.Wrap(err, "drop db error")
	}

	// Close the SQLite connection before dropping the database file
	err = m.DB().Close()
	if err != nil {
		m.Log().Errorf("drop dbPath error: %w", err) // NOTE: Maybe it was already closed.
	}

	err = os.Remove(m.dbPath)
	if err != nil {
		return m.dbPath, err
	}

	return m.dbPath, nil
}

func (m *Migrator) CloseAppConns() (string, error) {
	dbName := m.Cfg().ValOrDef("sql.database", "")

	// Close all open connections associated with the database
	err := m.DB().Close()
	if err != nil {
		return dbName, err
	}

	// Reopen the database connection for administrative tasks
	adminConn, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return dbName, err
	}
	defer adminConn.Close()

	// Terminate all connections to the database
	st := `SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = $1;`
	_, err = adminConn.Exec(st, dbName)
	if err != nil {
		return dbName, err
	}

	return dbName, nil
}

// DropDb migration.
func (m *Migrator) createMigrationsTable() (migTable string, err error) {
	tx, err := m.GetTx()
	if err != nil {
		return migTable, errors.Wrap(err, "create migration table error")
	}

	st := fmt.Sprintf(sqlCreateMigrationsSt, m.schema, sqlMigrationsTable)

	_, err = tx.Exec(st)
	if err != nil {
		return sqlMigrationsTable, err
	}

	return sqlMigrationsTable, tx.Commit()
}

func (m *Migrator) AddMigration(e Exec) {
	m.migs = append(m.migs, &Migration{Executor: e})
}

func (m *Migrator) Migrate() (err error) {
	err = m.PreSetup()
	if err != nil {
		return errors.Wrap(err, "migrate error")
	}

	for _, mg := range m.migs {
		exec := mg.Executor
		fn := getFxName(exec.GetUp())
		name := migName(fn)

		// Continue if already applied
		if !m.canApplyMigration(name) {
			m.Log().Infof("Migration '%s' already applied.", name)
			continue
		}

		// Get a new Tx from migrator
		tx, err := m.GetTx()
		if err != nil {
			return errors.Wrap(err, "migrate error")
		}

		// Pass Tx to the executor
		exec.SetTx(tx)

		// Execute migration
		values := reflect.ValueOf(exec).MethodByName(fn).Call([]reflect.Value{})

		// Read error
		err, ok := values[0].Interface().(error)
		if !ok && err != nil {
			m.Log().Infof("Migration not executed: %s", fn)   // TODO: Remove log
			m.Log().Infof("Err  %+v' of type %T\n", err, err) // TODO: Remove log.
			msg := fmt.Sprintf("migrate cannot run migration '%s': %s", fn, err.Error())
			err = tx.Rollback()
			if err != nil {
				return errors.Wrap(err, "migrate rollback error")
			}
			return errors.NewError(msg)
		}

		// Register migration
		err = m.recMigration(exec)

		err = tx.Commit()
		if err != nil {
			msg := fmt.Sprintf("Cannot update migrations table: %s\n", err.Error())
			m.Log().Errorf("migrate commit error: %s", msg)
			err = tx.Rollback()
			if err != nil {
				return errors.Wrap(err, "migrate rollback error")
			}
			return errors.NewError(msg)
		}

		log.Printf("Migration executed: %s\n", fn)
	}

	return nil
}

// Rollback migrations.
func (m *Migrator) Rollback(steps ...int) error {
	// Default to 1 step if no value is provided
	s := 1
	if len(steps) > 0 && steps[0] > 1 {
		s = steps[0]
	}

	// Default to max n° migration if steps is higher
	c := m.count()
	if s > c {
		s = c
	}

	m.rollback(s)
	return nil
}

// RollbackAll migrations.
func (m *Migrator) RollbackAll() error {
	return m.rollback(m.count())
}

func (m *Migrator) rollback(steps int) error {
	count := m.count()
	stopAt := count - steps

	for i := count - 1; i >= stopAt; i-- {
		mg := m.migs[i]
		exec := mg.Executor
		fn := getFxName(exec.GetDown())
		// Migration name is associated to up migration
		name := migName(getFxName(exec.GetUp()))

		// Continue if already not rolledback
		if m.cancelRollback(name) {
			log.Printf("Rollback '%s' already executed.", name)
			continue
		}

		// Get a new Tx from migrator
		tx, err := m.GetTx()
		if err != nil {
			return errors.Wrap(err, "rollback error")
		}

		// Pass Tx to the executor
		exec.SetTx(tx)

		// Execute rollback
		values := reflect.ValueOf(exec).MethodByName(fn).Call([]reflect.Value{})

		// Read error
		err, ok := values[0].Interface().(error)
		if !ok && err != nil {
			log.Printf("Rollback not executed: %s\n", fn)
			log.Printf("Err '%+v' of type %T", err, err)
		}

		// Remove migration record.
		err = m.delMigration(exec)

		err = tx.Commit()
		if err != nil {
			msg := fmt.Sprintf("Cannot update migrations table: %s\n", err.Error())
			log.Printf("Commit error: %s", msg)
			tx.Rollback()
			return errors.NewError(msg)
		}

		log.Printf("Rollback executed: %s\n", fn)
	}

	return nil
}

func (m *Migrator) SoftReset() error {
	err := m.RollbackAll()
	if err != nil {
		log.Printf("Cannot rollback database: %s", err.Error())
		return err
	}

	err = m.Migrate()
	if err != nil {
		log.Printf("Cannot migrate database: %s", err.Error())
		return err
	}

	return nil
}

func (m *Migrator) Reset() error {
	_, err := m.DropDb()
	if err != nil {
		m.Log().Errorf("Drop database error: %w", err)
		// Don't return maybe it was not created before.
	}

	_, err = m.CreateDb()
	if err != nil {
		return errors.Wrap(err, "create database error")
	}

	err = m.Migrate()
	if err != nil {
		return errors.Wrap(err, "drop database error")
	}

	return nil
}

func (m *Migrator) recMigration(e Exec) error {
	st := fmt.Sprintf(sqlRecMigrationSt, m.schema, sqlMigrationsTable)
	upFx := getFxName(e.GetUp())
	downFx := getFxName(e.GetDown())
	name := migName(upFx)
	log.Printf("%+s", upFx)

	uid, err := uuid.New()
	if err != nil {
		return errors.Wrap(err, "rec migration error")
	}

	_, err = m.DB().Exec(st,
		uid,
		ToNullString(name),
		ToNullString(upFx),
		ToNullString(downFx),
		ToNullBool(true),
		ToNullTime(time.Time{}),
	)

	if err != nil {
		return errors.Wrap(err, "cannot update migrations table")
	}

	return nil
}

func (m *Migrator) cancelRollback(name string) bool {
	st := fmt.Sprintf(sqlSelMigrationSt, m.schema, sqlMigrationsTable, name)
	r, err := m.DB().Query(st)

	if err != nil {
		m.Log().Errorf("Cannot determine rollback status: %w", err)
		return true
	}

	for r.Next() {
		var applied sql.NullBool
		err = r.Scan(&applied)
		if err != nil {
			m.Log().Errorf("Cannot determine migration status: %w", err)
			return true
		}

		return !applied.Bool
	}

	return true
}

func (m *Migrator) canApplyMigration(name string) bool {
	st := fmt.Sprintf(sqlSelMigrationSt, m.schema, sqlMigrationsTable, name)
	r, err := m.DB().Query(st)

	if err != nil {
		m.Log().Errorf("Cannot determine migration status: %w", err)
		return false
	}

	for r.Next() {
		var applied sql.NullBool
		err = r.Scan(&applied)
		if err != nil {
			m.Log().Errorf("Cannot determine migration status: %s", err)
			return false
		}

		return !applied.Bool
	}

	return true
}

func (m *Migrator) delMigration(e Exec) error {
	name := migName(getFxName(e.GetUp()))
	st := fmt.Sprintf(sqlDelMigrationSt, m.schema, sqlMigrationsTable, name)
	_, err := e.GetTx().Exec(st)

	if err != nil {
		return errors.Wrap(err, "cannot update migrations table")
	}

	return nil
}

func (m *Migrator) count() (last int) {
	return len(m.migs)
}

func (m *Migrator) last() (last int) {
	return m.count() - 1
}

func getFxName(i interface{}) string {
	n := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	t := strings.FieldsFunc(n, split)
	return t[len(t)-2]
}

func split(r rune) bool {
	return r == '.' || r == '-'
}

func migName(upFxName string) string {
	return toSnakeCase(upFxName)
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

type NullTime struct {
	Time  time.Time
	Valid bool // Indicates if the timestamp is null or not
}

func ToNullTime(t time.Time) NullTime {
	return NullTime{
		Time:  t,
		Valid: true,
	}
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func ToNullBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}
