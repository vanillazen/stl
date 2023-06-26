package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver import

	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/config"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
)

type DB struct {
	sys.Core
	db *sql.DB
}

var (
	cfgKey = config.Key
)

func NewDB(opts ...sys.Option) *DB {
	return &DB{
		Core: sys.NewCore("sqlite-db", opts...),
	}
}

func (db *DB) DB() *sql.DB {
	return db.db
}

func (db *DB) Start(ctx context.Context) error {
	return db.Connect(ctx)
}

func (db *DB) Connect(ctx context.Context) error {
	sqlDB, err := sql.Open("sqlite3", db.dbPath())
	if err != nil {
		msg := fmt.Sprintf("%s connection error", db.Name())
		return errors.Wrap(err, msg)
	}

	err = sqlDB.Ping()
	if err != nil {
		msg := fmt.Sprintf("%s ping connection error", db.Name())
		return errors.Wrap(err, msg)
	}

	db.db = sqlDB
	db.Log().Infof("%s database connected", db.Name())
	return nil
}

func (db *DB) DBConn(ctx context.Context) (*sql.DB, error) {
	return db.db, nil
}

func (db *DB) dbPath() string {
	cfg := db.Cfg()
	dbPath := cfg.GetString(cfgKey.SQLiteFilePath)
	return dbPath
}
