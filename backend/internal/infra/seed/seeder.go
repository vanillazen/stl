package seed

import (
	"github.com/vanillazen/stl/backend/internal/sys"
)

type (
	Seeder interface {
		sys.Core
		// Seed applies pending seeding
		Seed() error
		// SetAssetsPath sets the path form where the seeding are read
		SetAssetsPath(path string)
		// AssetsPath returns the path form where the seeding are read
		AssetsPath() string
	}
)
