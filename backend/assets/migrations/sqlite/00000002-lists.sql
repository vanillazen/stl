--UP
CREATE TABLE lists (
                       id TEXT PRIMARY KEY,
                       name TEXT NOT NULL,
                       description TEXT NOT NULL,
                       owner_id TEXT NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

--DOWN
DROP TABLE lists;
