package sqlite

const (
	createSeederTable = `CREATE TABLE %s (
    id UUID PRIMARY KEY,
    idx INTEGER,
    name VARCHAR(64),
    created_at TEXT
    );`

	droopSeederTable = `DROP TABLE %s;`

	selFromSeedsTable = `SELECT (COUNT(*) > 0) AS record_exists 
		FROM %s 
		WHERE idx = %d 
		    AND name = '%s'`

	insertSeederTable = `INSERT INTO %s (id, idx, name, created_at)
	VALUES (:id, :idx, :name, :created_at);`

	deleteSeederTable = `DELETE FROM %s WHERE idx = %d AND name = '%s'`
)
