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
		up_fx VARCHAR(64),
		down_fx VARCHAR(64),
 		is_applied BOOLEAN,
		created_at TIMESTAMP
	);`

	sqlDropMigrationsSt = `DROP TABLE %s;`

	selMigrationSt = `
	SELECT is_applied FROM %s 
	                  WHERE name = '%s' and is_applied = true`

	recMigrationSt = `
	INSERT INTO %s (id, name, up_fx, down_fx, is_applied, created_at)
	VALUES (:id, :name, :up_fx, :down_fx, :is_applied, :created_at);`

	delMigrationSt = `DELETE FROM %s 
       WHERE name = '%s' 
           AND is_applied = true`
)
