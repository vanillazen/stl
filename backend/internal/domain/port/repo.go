package port

import (
	"context"

	"github.com/vanillazen/stl/backend/internal/domain/model"
	"github.com/vanillazen/stl/backend/internal/infra/db"
	"github.com/vanillazen/stl/backend/internal/sys"
)

type (
	Repo interface {
		sys.Core
		DB(ctx context.Context) (db db.DB)
	}

	ListRepo interface {
		Repo
		// CreateList in persistence
		CreateList(ctx context.Context, list model.List) (model.List, error)
		// GetList from persistence
		GetList(ctx context.Context, userID string, preload ...bool) (list model.List, err error)
		//// UpdateList in persistence
		//UpdateList(ctx context.Context, task model.List, userID string) error
		//// DeleteList in persistence
		//DeleteList(ctx context.Context, listID, userID string) error
		//
		//// AddTask in persistence
		//AddTask(ctx context.Context, list string, task model.Task, userID string) (model.Task, error)
		//// AddTasks in persistence
		//AddTasks(ctx context.Context, list string, task []model.Task, userID string) error
		//// GetTask from persistence
		//GetTask(ctx context.Context, taskID, userID string) (user model.Task, err error)
		//// UpdateTask in persistence
		//UpdateTask(ctx context.Context, task *model.Task, userID string) error
		//// DeleteTask in persistence
		//DeleteTask(ctx context.Context, taskID, userID string) error
		//
		// GetUser from persistence
		GetUser(ctx context.Context, userID string) (user model.User, err error)
	}
)
