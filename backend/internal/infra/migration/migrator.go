package migrator

import (
	"github.com/vanillazen/stl/backend/internal/sys"
)

type (
	Migrator interface {
		sys.Core
		// Migrate applies pending migrations
		Migrate() error
		// Rollback reverts from one to N migrations already applied
		Rollback(steps ...int) error
		// RollbackAll reverts all migrations allready applied
		RollbackAll() error
		// SoftReset apply all migrations again after rolling back all migrations.
		SoftReset() error
		// Reset apply all migrations again after dropping the database and recreating it
		Reset() error
	}
)
