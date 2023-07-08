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

func (r *ListRepo) GetList(ctx context.Context, userID string, preload ...bool) (list model.List, err error) {
	dbase := r.DB(ctx).DB()

	query := `
		SELECT l.id, l.name, l.description, u.id, u.username, u.name, u.email, u.password, l.created_at, l.updated_at,
			t.id, t.list_id, t.name, t.description, t.category, t.tags, t.location, t.created_at, t.updated_at
		FROM lists l
		INNER JOIN users u ON l.owner_id = u.id
		LEFT JOIN tasks t ON l.id = t.list_id
		WHERE u.id = $1
		ORDER BY l.id, t.id
	`

	rows, err := dbase.Query(query, userID)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	var currentListID model.ID
	var currentTaskID model.ID

	for rows.Next() {
		var task model.Task

		err := rows.Scan(
			list.ID.String(),
			list.Name,
			list.Description,
			list.CreatedAt,
			list.UpdatedAt,
			task.ID.String(),
			task.ListID.String(),
			task.Name,
			task.Description,
			task.Category,
			task.Tags,
			task.Location,
			task.CreatedAt,
			task.UpdatedAt,
		)
		if err != nil {
			return list, err
		}

		if currentListID.String() != list.ID.String() {
			currentListID = list.ID
			list.Tasks = []model.Task{task}
		} else {
			if currentTaskID.String() != task.ID.String() {
				currentTaskID = task.ID
				list.Tasks = append(list.Tasks, task)
			}
		}
	}

	return list, nil
}

func (r *ListRepo) GetUser(ctx context.Context, userID string) (user model.User, err error) {
	// WIP: Mock implementation
	ref := "e1263c73-521b-41b5-96e5-58c3f71e65a1\""

	ok := uuid.Validate(userID)
	if !ok {
		return user, InvalidResourceIDErr
	}

	uid, err := uuid.Parse(ref)
	if err != nil {
		return user, err
	}

	if userID == ref {
		return model.User{
			ID:       model.NewID(uid),
			Username: "johndoe",
			Name:     "John Doe",
			Email:    "john.doe@localhost.com",
		}, nil
	}

	return model.User{}, UserNotFoundErr
}
