package sqlite

const (
	migrationsTable = "migration"

	createDBSt = `
		CREATE DATABASE %s;`

	dropDBSt = `
		DROP DATABASE %s;`

	createMigraationTableSt = `CREATE TABLE %s (
    id UUID PRIMARY KEY,
		name VARCHAR(64),
		created_at TEXT
	);`

	sqlDropMigrationsSt = `DROP TABLE %s;`

	selMigrationSt = `SELECT (COUNT(*) > 0) AS record_exists FROM %s WHERE name = '%s'`

	recMigrationSt = `
	INSERT INTO %s (id, name, created_at)
	VALUES (:id, :name, :created_at);`

	delMigrationSt = `DELETE FROM %s 
       WHERE name = '%s'`
)
