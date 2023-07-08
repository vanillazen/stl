package fixture

import "github.com/vanillazen/stl/backend/internal/sys"

type (
	Fixture interface {
		sys.Core
		PopulateDB() error
	}
)
