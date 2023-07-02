--UP
CREATE TABLE tasks (
                       id TEXT PRIMARY KEY,
                       list_id TEXT NOT NULL,
                       name TEXT NOT NULL,
                       description TEXT NOT NULL,
                       category TEXT[],
                       tags TEXT[],
                       location TEXT[],
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (list_id) REFERENCES lists (id) ON DELETE CASCADE
);

--DOWN
DROP TABLE tasks;
