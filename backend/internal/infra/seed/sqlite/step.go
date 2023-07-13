package sqlite

import "database/sql"

type (
	step struct {
		Index int64
		Name  string
		Seeds []SeedFx
		tx    *sql.Tx
	}
)

func (s *step) Config(seed []SeedFx) {
	s.Seeds = seed
}

func (s *step) GetIndex() (idx int64) {
	return s.Index
}

func (s *step) GetName() (name string) {
	return s.Name
}

func (s *step) GetSeeds() (seeds []SeedFx) {
	return s.Seeds
}

func (s *step) SetTx(tx *sql.Tx) {
	s.tx = tx
}

func (s *step) GetTx() (tx *sql.Tx) {
	return s.tx
}
