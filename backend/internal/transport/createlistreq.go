package transport

import "github.com/vanillazen/stl/backend/internal/domain/model"

type (
	CreateListReq struct {
		UserID      string
		Name        string
		Description string
	}
)

func (req CreateListReq) ToList() model.List {
	return model.List{
		Name:        req.Name,
		Description: req.Description,
	}
}
