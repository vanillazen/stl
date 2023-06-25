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
		Connect(ctx context.Context) error
	}
)
