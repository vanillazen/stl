package model

type (
	Task struct {
		ID
		ListID      ID
		Name        string
		Description string
		Category    []string
		Tags        []string
		Location    []string
		Audit
	}
)
