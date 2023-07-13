package sqlite

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/vanillazen/stl/backend/internal/domain/model"
	"github.com/vanillazen/stl/backend/internal/infra/db"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/uuid"
)

type (
	// Fixture struct.
	Fixture struct {
		sys.Core
		db db.DB
	}
)

func NewFixture(db db.DB, opts ...sys.Option) (fxt *Fixture) {
	f := &Fixture{
		Core: sys.NewCore("fixture", opts...),
		db:   db,
	}

	return f
}

type (
	User struct {
		ID        string
		Username  string
		Name      string
		Email     string
		Password  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	List struct {
		ID          string
		Name        string
		Description string
		OwnerID     string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	Task struct {
		ID          string
		ListID      string
		Name        string
		Description string
		Category    []string
		Tags        []string
		Location    []string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)

func (f *Fixture) DB() *sql.DB {
	return f.db.DB()
}

func (f *Fixture) Start(ctx context.Context) error {
	return f.PopulateDB()
}

func (f *Fixture) PopulateDB() error {
	dbase := f.DB()

	// Users
	users := []model.User{
		{
			ID:       model.NewID(uuid.MustParse("0792b97b-4f88-42a8-a035-1d0aad0ae7f8")),
			Username: "user1",
			Name:     "User 1",
			Email:    "user1@example.com",
			Password: "password1",
		},
		{
			ID:       model.NewID(uuid.NewUUID()),
			Username: "user2",
			Name:     "User 2",
			Email:    "user2@example.com",
			Password: "password2",
		},
		{
			ID:       model.NewID(uuid.NewUUID()),
			Username: "user3",
			Name:     "User 3",
			Email:    "user3@example.com",
			Password: "password3",
		},
	}

	for _, user := range users {
		insertUser(dbase, user)
	}

	// Lists
	lists := []model.List{
		{
			ID:          model.NewID(uuid.MustParse("cdc7a443-3c6a-431b-b45a-b14735953a19")),
			Name:        "List 1",
			Description: "List 1 Description",
			Owner:       users[0],
			Tasks:       []model.Task{},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
		{
			ID:          model.NewID(uuid.NewUUID()),
			Name:        "List 2",
			Description: "List 2 Description",
			Owner:       users[1],
			Tasks:       []model.Task{},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
		{
			ID:          model.NewID(uuid.NewUUID()),
			Name:        "List 3",
			Description: "List 3 Description",
			Owner:       users[1],
			Tasks:       []model.Task{},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
		{
			ID:          model.NewID(uuid.NewUUID()),
			Name:        "List 4",
			Description: "List 4 Description",
			Owner:       users[1],
			Tasks:       []model.Task{},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
		{
			ID:          model.NewID(uuid.NewUUID()),
			Name:        "List 5",
			Description: "List 5 Description",
			Owner:       users[2],
			Tasks:       []model.Task{},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
	}

	for _, list := range lists {
		insertList(dbase, list)
	}

	// Tasks
	tasks := []model.Task{
		{
			ID:          model.NewID(uuid.NewUUID()),
			ListID:      model.NewID(lists[0].ID.UUID),
			Name:        "Task 1",
			Description: "Task 1 Description",
			Category:    []string{"Category 1"},
			Tags:        []string{"Tag 1", "Tag 2"},
			Location:    []string{"Location 1"},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
		{
			ID:          model.NewID(uuid.NewUUID()),
			ListID:      model.NewID(lists[0].ID.UUID),
			Name:        "Task 2",
			Description: "Task 2 Description",
			Category:    []string{"Category 2"},
			Tags:        []string{"Tag 3", "Tag 4"},
			Location:    []string{"Location 2"},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
		{
			ID:          model.NewID(uuid.NewUUID()),
			ListID:      model.NewID(lists[1].ID.UUID),
			Name:        "Task 3",
			Description: "Task 3 Description",
			Category:    []string{"Category 3"},
			Tags:        []string{"Tag 5", "Tag 6"},
			Location:    []string{"Location 3"},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
		{
			ID:          model.NewID(uuid.NewUUID()),
			ListID:      model.NewID(lists[2].ID.UUID),
			Name:        "Task 4",
			Description: "Task 4 Description",
			Category:    []string{"Category 4"},
			Tags:        []string{"Tag 7", "Tag 8"},
			Location:    []string{"Location 4"},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
		{
			ID:          model.NewID(uuid.NewUUID()),
			ListID:      model.NewID(lists[3].ID.UUID),
			Name:        "Task 5",
			Description: "Task 5 Description",
			Category:    []string{"Category 5"},
			Tags:        []string{"Tag 9", "Tag 10"},
			Location:    []string{"Location 5"},
			Audit:       model.NewAudit(time.Now(), time.Now()),
		},
	}

	for _, task := range tasks {
		insertTask(dbase, task)
	}

	return nil
}

func insertUser(db *sql.DB, user model.User) {
	query := `
		INSERT INTO users (id, username, name, email, password)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, user.ID.UUID.Val, user.Username, user.Name, user.Email, user.Password)
	if err != nil {
		log.Fatal(err)
	}
}

func insertList(db *sql.DB, list model.List) {
	query := `
		INSERT INTO lists (id, name, description, owner_id)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(query, list.ID.UUID.Val, list.Name, list.Description, list.Owner.ID.UUID.Val)
	if err != nil {
		log.Fatal(err)
	}
}

func insertTask(db *sql.DB, task model.Task) {
	query := `
		INSERT INTO tasks (id, list_id, name, description, category, tags, location)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, task.ID.UUID.Val, task.ListID.UUID.Val, task.Name, task.Description, task.Category, task.Tags, task.Location)
	if err != nil {
		log.Fatal(err)
	}
}
