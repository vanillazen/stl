package model

type (
	List struct {
		ID
		Name        string
		Description string
		Owner       User
		Tasks       []*Task
		Audit
	}
)
