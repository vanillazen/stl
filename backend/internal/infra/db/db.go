package db

import (
	"context"
	"database/sql"

	"github.com/vanillazen/stl/backend/internal/sys"
)

type (
	DB interface {
		sys.Core
		DB() *sql.DB
		Path() string
		Schema() string
		Name() string
		Connect(ctx context.Context) error
	}
)
