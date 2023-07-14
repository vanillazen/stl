package sqlite

const (
	createSeedsTable = `CREATE TABLE %s (
    id UUID PRIMARY KEY,
    idx INTEGER,
    name VARCHAR(64),
    created_at TEXT
    );`

	dropSeedsTable = `DROP TABLE %s;`

	selectFromSeeds = `SELECT (COUNT(*) > 0) AS record_exists 
		FROM %s 
		WHERE idx = %d 
		    AND name = '%s'`

	insertIntoSeeds = `INSERT INTO %s (id, idx, name, created_at)
	VALUES (:id, :idx, :name, :created_at);`

	deleteFromSeeds = `DELETE FROM %s WHERE idx = %d AND name = '%s'`
)
