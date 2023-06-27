package sqlite

import (
	"database/sql"
	fmt "fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/vanillazen/stl/backend/internal/sys/errors"
)

const (
	migPath = "assets/migrations/sqlite"
)

func (m *Migrator) addSteps() error {
	qq, err := m.readMigQueries()
	if err != nil {
		return err
	}

	for i, q := range qq {
		s := &step{
			Name: q.Name,
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(q.Up)
				m.Log().Error(err)
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec(q.Down)
				return err
			},
		}

		m.AddMigration(i, s)
	}

	return nil
}

type queries struct {
	Name string
	Up   string
	Down string
}

func (m *Migrator) readMigQueries() ([]queries, error) {
	var qq []queries

	files, err := m.fs.ReadDir(migPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", migPath, file.Name())
		content, err := m.fs.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		sections := strings.Split(string(content), "--DOWN")
		if len(sections) < 2 {
			msg := fmt.Sprintf("invalid migration file format: %s", file.Name())
			return nil, errors.NewError(msg)
		}

		up := strings.TrimSpace(strings.TrimPrefix(sections[0], "--UP\n"))
		down := strings.TrimSpace(sections[1])

		q := queries{
			Name: stepName(filePath),
			Up:   up,
			Down: down,
		}

		qq = append(qq, q)
	}

	return qq, nil
}

func stepName(filename string) string {
	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, filepath.Ext(base))

	re := regexp.MustCompile(`^[-\d]+`)
	name := re.ReplaceAllString(base, "")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	return name
}
