package transport

import (
	"time"

	"github.com/vanillazen/stl/backend/internal/domain/model"
	"github.com/vanillazen/stl/backend/internal/sys/config"
	v "github.com/vanillazen/stl/backend/internal/sys/validator"
)

type (
	CreateListRes struct {
		ServiceRes
		UserID      string
		Name        string
		Description string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)

func NewCreateListRes(valErrSet v.ValErrorSet, err error, cfg *config.Config) CreateListRes {
	return CreateListRes{
		ServiceRes: NewServiceRes(valErrSet, err, cfg),
	}
}

func (res *CreateListRes) FromList(m model.List) {
	res.UserID = m.Owner.ID.String()
	res.Name = m.Name
	res.Description = m.Description
	res.CreatedAt = m.CreatedAt
	res.UpdatedAt = m.UpdatedAt
}
