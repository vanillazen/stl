package sqlite

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

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
	users := []User{
		{
			ID:        "0792b97b-4f88-42a8-a035-1d0aad0ae7f8",
			Username:  "user1",
			Name:      "User 1",
			Email:     "user1@example.com",
			Password:  "password1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        genID(),
			Username:  "user2",
			Name:      "User 2",
			Email:     "user2@example.com",
			Password:  "password2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        genID(),
			Username:  "user3",
			Name:      "User 3",
			Email:     "user3@example.com",
			Password:  "password3",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		insertUser(dbase, user)
	}

	// Lists
	lists := []List{
		{
			ID:          "cdc7a443-3c6a-431b-b45a-b14735953a19",
			Name:        "List 1",
			Description: "List 1 Description",
			OwnerID:     users[0].ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			Name:        "List 2",
			Description: "List 2 Description",
			OwnerID:     users[1].ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			Name:        "List 3",
			Description: "List 3 Description",
			OwnerID:     users[1].ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			Name:        "List 4",
			Description: "List 4 Description",
			OwnerID:     users[1].ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			Name:        "List 5",
			Description: "List 5 Description",
			OwnerID:     users[2].ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, list := range lists {
		insertList(dbase, list)
	}

	tasks := []Task{
		{
			ID:          genID(),
			ListID:      lists[0].ID,
			Name:        "Task 1",
			Description: "Task 1 Description",
			Category:    []string{"Category 1"},
			Tags:        []string{"Tag 1", "Tag 2"},
			Location:    []string{"Location 1"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			ListID:      lists[0].ID,
			Name:        "Task 2",
			Description: "Task 2 Description",
			Category:    []string{"Category 2"},
			Tags:        []string{"Tag 2", "Tag 3"},
			Location:    []string{"Location 2"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			ListID:      lists[0].ID,
			Name:        "Task 3",
			Description: "Task 3 Description",
			Category:    []string{"Category 1", "Category 2"},
			Tags:        []string{"Tag 1", "Tag 3"},
			Location:    []string{"Location 1", "Location 2"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			ListID:      lists[1].ID,
			Name:        "Task 4",
			Description: "Task 4 Description",
			Category:    []string{"Category 3"},
			Tags:        []string{"Tag 3", "Tag 4"},
			Location:    []string{"Location 3"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			ListID:      lists[1].ID,
			Name:        "Task 5",
			Description: "Task 5 Description",
			Category:    []string{"Category 1", "Category 3"},
			Tags:        []string{"Tag 1", "Tag 4"},
			Location:    []string{"Location 1", "Location 3"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			ListID:      lists[2].ID,
			Name:        "Task 6",
			Description: "Task 6 Description",
			Category:    []string{"Category 2", "Category 3"},
			Tags:        []string{"Tag 2", "Tag 4"},
			Location:    []string{"Location 2", "Location 3"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          genID(),
			ListID:      lists[3].ID,
			Name:        "Task 7",
			Description: "Task 7 Description",
			Category:    []string{"Category 1", "Category 3"},
			Tags:        []string{"Tag 1", "Tag 4"},
			Location:    []string{"Location 1", "Location 3"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, task := range tasks {
		insertTask(dbase, task)
	}

	return nil
}

func genID() string {
	return uuid.Must().String()
}

func insertUser(db *sql.DB, user User) {
	query := `
		INSERT INTO users (id, username, name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, user.ID, user.Username, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Fatal(err)
	}
}

func insertList(db *sql.DB, list List) {
	query := `
		INSERT INTO lists (id, name, description, owner_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, list.ID, list.Name, list.Description, list.OwnerID, list.CreatedAt, list.UpdatedAt)
	if err != nil {
		log.Fatal(err)
	}
}

func insertTask(db *sql.DB, task Task) {
	query := `
		INSERT INTO tasks (id, list_id, name, description, category, tags, location, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	categoryStr := strings.Join(task.Category, ",")
	tagsStr := strings.Join(task.Tags, ",")
	locationStr := strings.Join(task.Location, ",")

	_, err := db.Exec(query, task.ID, task.ListID, task.Name, task.Description, categoryStr, tagsStr, locationStr, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		log.Fatal(err)
	}
}
