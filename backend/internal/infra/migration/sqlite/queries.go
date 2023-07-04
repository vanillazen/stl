package sqlite

const (
	createDB = `CREATE DATABASE %s;`

	dropDB = `DROP DATABASE %s;`

	createMigTable = `CREATE TABLE %s (
    id UUID PRIMARY KEY,
    idx INTEGER,
    name VARCHAR(64),
    created_at TEXT
    );`

	dropMigTable = `DROP TABLE %s;`

	selFromMigTable = `SELECT (COUNT(*) > 0) AS record_exists 
		FROM %s 
		WHERE idx = %d 
		    AND name = '%s'`

	insertMigTable = `INSERT INTO %s (id, idx, name, created_at)
	VALUES (:id, :idx, :name, :created_at);`

	delFromMigTable = `DELETE FROM %s WHERE idx = %d AND name = '%s'`
)
