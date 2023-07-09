package model

type (
	Task struct {
		ID
		ListID      ID
		Name        string
		Description string
		Category    StringSlice
		Tags        StringSlice
		Location    StringSlice
		Audit
	}
)
