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

func NewListRepo(db db.DB, opts ...sys.Option) (*ListRepo, error) {
	return &ListRepo{
		SimpleCore: sys.NewCore("list-repo", opts...),
		db:         db,
	}, nil
}

func (cr *ListRepo) Setup(ctx context.Context) error {
	err := cr.db.Connect(ctx)
	if err != nil {
		err = errors.Wrap("list repo setup error", err)
		return err
	}

	err = cr.db.Connect(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ListRepo) DB(ctx context.Context) db.DB {
	return cr.db
}

func (cr *ListRepo) CreateList(ctx context.Context, m model.List) (model.List, error) {
	return m, errors.NewError("not implemented yet")
}

func (cr *ListRepo) GetUser(ctx context.Context, userID string) (model.User, error) {
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
