package service

import (
	"context"
	"fmt"

	"github.com/vanillazen/stl/backend/internal/domain/port"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
	t "github.com/vanillazen/stl/backend/internal/transport"
)

type (
	ListService interface {
		sys.Core
		CreateList(ctx context.Context, req t.CreateListReq) t.CreateListRes
		//GetList(...)
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
		err = errors.Wrap("create list error", err)
		return t.NewCreateListRes(nil, err, rs.Cfg())
	}

	list.Owner = user

	// Persist it
	_, err = rs.Repo().CreateList(ctx, list)
	if err != nil {
		err = errors.Wrap("create list error", err)
		return t.NewCreateListRes(nil, err, rs.Cfg())
	}

	return t.NewCreateListRes(nil, nil, nil)
}

func (rs *List) Repo() port.ListRepo {
	return rs.repo
}

func (rs *List) Start(ctx context.Context) error {
	db := rs.repo.DB(ctx)

	err := db.Start(ctx)
	if err != nil {
		msg := fmt.Sprintf("%s start error", rs.Name())
		return errors.Wrap(msg, err)
	}

	return nil
}
