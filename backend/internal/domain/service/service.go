package service

import (
	"context"

	"github.com/vanillazen/stl/backend/internal/domain/port"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
	t "github.com/vanillazen/stl/backend/internal/transport"
)

type (
	ListService interface {
		sys.Core
		CreateList(ctx context.Context, req t.CreateListReq) t.CreateListRes
		GetList(ctx context.Context, req t.GetListReq) t.GetListRes
		//UpdateList(...)
		//DeleteList(...)
		//AddTask(...)
		//AddTasks(...)
		//GetTask(...)
		//UpdateTask(...)
		//DeleteTask(...)
		//GetUser(...)
	}

	List struct {
		*sys.SimpleCore
		repo   port.ListRepo
		mailer port.Mailer
	}
)

func NewService(rr port.ListRepo, opts ...sys.Option) *List {
	return &List{
		SimpleCore: sys.NewCore("list-service", opts...),
		repo:       rr,
		mailer:     nil, // Interface not implemented yet
	}
}

func (rs *List) CreateList(ctx context.Context, req t.CreateListReq) (res t.CreateListRes) {
	// Transport to Model
	list := req.ToList()

	// Validate model
	v := NewListValidator(list)

	err := v.ValidateForCreate()
	if err != nil {
		return t.NewCreateListRes(v.Errors, err, rs.Cfg())
	}

	// Set Owner
	user, err := rs.Repo().GetUser(ctx, req.UserID)
	if err != nil {
		err = errors.Wrap(err, "create list error")
		return t.NewCreateListRes(nil, err, rs.Cfg())
	}

	list.Owner = user

	// Persist it
	_, err = rs.Repo().CreateList(ctx, list)
	if err != nil {
		err = errors.Wrap(err, "create list error")
		return t.NewCreateListRes(nil, err, rs.Cfg())
	}

	return t.NewCreateListRes(nil, nil, nil)
}

func (rs *List) GetList(ctx context.Context, req t.GetListReq) (res t.GetListRes) {
	_, err := rs.Repo().GetList(ctx, req.UserID, true)
	if err != nil {
		err = errors.Wrap(err, "get list repo error")
		return t.NewGetListRes(nil, err, rs.Cfg())
	}

	return t.NewGetListRes(nil, nil, nil)
}

func (rs *List) Repo() port.ListRepo {
	return rs.repo
}
