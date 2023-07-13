package migrator

import (
	"github.com/vanillazen/stl/backend/internal/sys"
)

type (
	Migrator interface {
		sys.Core
		// Migrate applies pending seeding
		Seed() error
		// Rollback reverts from one to N seeding already applied
		Rollback(steps ...int) error
		// RollbackAll reverts all seeding allready applied
		RollbackAll() error
		// SoftReset apply all seeding again after rolling back all seeding.
		SoftReset() error
		// Reset apply all seeding again after dropping the database and recreating it
		Reset() error
		// SetAssetsPath sets the path form where the seeding are read
		SetAssetsPath(path string)
		// AssetsPath returns the path form where the seeding are read
		AssetsPath() string
	}
)
