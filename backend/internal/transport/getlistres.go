package transport

import (
	"time"

	"github.com/vanillazen/stl/backend/internal/domain/model"
	"github.com/vanillazen/stl/backend/internal/sys/config"
	v "github.com/vanillazen/stl/backend/internal/sys/validator"
)

type (
	GetListRes struct {
		ServiceRes
		UserID      string
		Name        string
		Description string
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Tasks       []Task
	}

	Task struct {
		Name        string
		Description string
		Category    []string
		Tags        []string
		Location    []string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)

func NewGetListRes(valErrSet v.ValErrorSet, err error, cfg *config.Config, list model.List) GetListRes {
	var tasks []Task
	for _, m := range list.Tasks {
		tasks = append(tasks, Task{
			Name:        m.Name,
			Description: m.Description,
			Category:    m.Category,
			Tags:        m.Tags,
			Location:    m.Location,
		})
	}

	return GetListRes{
		ServiceRes:  NewServiceRes(valErrSet, err, cfg),
		Name:        list.Name,
		Description: list.Description,
		Tasks:       tasks,
	}
}

func (res *GetListRes) FromList(m model.List) {
	res.UserID = m.Owner.ID.String()
	res.Name = m.Name
	res.Description = m.Description
	res.CreatedAt = m.CreatedAt
	res.UpdatedAt = m.UpdatedAt
	for _, t := range m.Tasks {
		res.Tasks = append(res.Tasks,
			Task{
				Name:        t.Name,
				Description: t.Description,
				Category:    t.Category,
				Tags:        t.Tags,
				Location:    t.Location,
				CreatedAt:   t.CreatedAt,
				UpdatedAt:   t.UpdatedAt,
			},
		)
	}
}
