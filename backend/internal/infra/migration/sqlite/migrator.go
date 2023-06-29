package sqlite

import (
	"database/sql"
	fmt "fmt"
	"path/filepath"
	"regexp"
	"strconv"
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
			Index: q.Index,
			Name:  q.Name,
			Up:    m.genTxExecFunc(q.Up),
			Down:  m.genTxExecFunc(q.Down),
		}

		m.AddMigration(i, s)
	}

	return nil
}

func (m *Migrator) genTxExecFunc(query string) func(tx *sql.Tx) error {
	return func(tx *sql.Tx) error {
		_, err := tx.Exec(query)
		//m.Log().Debugf("%s", err.Error())
		return err
	}
}

type queries struct {
	Index int64
	Name  string
	Up    string
	Down  string
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

		idx, name := stepName(filePath)

		q := queries{
			Index: idx,
			Name:  name,
			Up:    up,
			Down:  down,
		}

		qq = append(qq, q)
	}

	return qq, nil
}

func stepName(filename string) (idx int64, name string) {
	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, filepath.Ext(base))

	re := regexp.MustCompile(`^[-\d]+`)
	indexStr := re.FindString(base)
	idx, _ = strconv.ParseInt(strings.TrimSuffix(indexStr, "-"), 10, 64)

	name = re.ReplaceAllString(base, "")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	return idx, name
}
