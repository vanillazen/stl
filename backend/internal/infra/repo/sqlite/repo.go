package sqlite

import (
	"context"

	"github.com/vanillazen/stl/backend/internal/domain/model"
	"github.com/vanillazen/stl/backend/internal/infra/db"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
	"github.com/vanillazen/stl/backend/internal/sys/uuid"
)

type ListRepo struct {
	*sys.SimpleCore
	db db.DB
}

func NewListRepo(db db.DB, opts ...sys.Option) *ListRepo {
	return &ListRepo{
		SimpleCore: sys.NewCore("list-repo", opts...),
		db:         db,
	}
}

func (r *ListRepo) DB(ctx context.Context) db.DB {
	return r.db
}

func (r *ListRepo) Setup(ctx context.Context) error {
	err := r.db.Connect(ctx)
	if err != nil {
		err = errors.Wrap(err, "list repo setup error")
		return err
	}

	return nil
}

func (r *ListRepo) Start(ctx context.Context) error {
	r.Log().Infof("%s started", r.Name())
	return nil
}

func (r *ListRepo) Stop(ctx context.Context) error {
	err := r.DB(ctx).DB().Close()
	if err != nil {
		err := errors.Wrapf(err, "%s stop error", r.Name())
		return err
	}

	r.Log().Infof("%s stopped", r.Name())
	return nil
}

func (r *ListRepo) CreateList(ctx context.Context, m model.List) (updated model.List, err error) {
	err = m.GenID()
	if err != nil {
		return m, errors.Wrap(err, "create list repo error")
	}

	//db := r.db.DB()
	//if err != nil {
	//	return m, errors.Wrap(err, "create list repo err")
	//}
	//
	return updated, nil
}

func (r *ListRepo) GetUser(ctx context.Context, userID string) (model.User, error) {
	// WIP: Mock implementation
	ref := "e1263c73-521b-41b5-96e5-58c3f71e65a1\""

	ok := uuid.Validate(userID)
	if !ok {
		return model.User{}, InvalidResourceIDErr
	}

	if userID == ref {
		return model.User{
			ID:       model.NewID(uuid.NewUUID(ref)),
			Username: "johndoe",
			Name:     "John Doe",
			Email:    "john.doe@localhost.com",
		}, nil
	}

	return model.User{}, UserNotFoundErr
}
