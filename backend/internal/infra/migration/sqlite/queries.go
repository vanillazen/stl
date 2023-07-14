package sqlite

const (
	createMigrationsDB = `CREATE DATABASE %s;`

	dropMigrationsDB = `DROP DATABASE %s;`

	createMigrationsTable = `CREATE TABLE %s (
    id UUID PRIMARY KEY,
    idx INTEGER,
    name VARCHAR(64),
    created_at TEXT
    );`

	dropMigrationsTable = `DROP TABLE %s;`

	selectFromMigrations = `SELECT (COUNT(*) > 0) AS record_exists 
		FROM %s 
		WHERE idx = %d 
		    AND name = '%s'`

	insertIntoMigrations = `INSERT INTO %s (id, idx, name, created_at)
	VALUES (:id, :idx, :name, :created_at);`

	deleteFromMigrations = `DELETE FROM %s WHERE idx = %d AND name = '%s'`
)
