--UP
CREATE TABLE users (
                       id TEXT PRIMARY KEY,
                       username TEXT NOT NULL,
                       name TEXT NOT NULL,
                       email TEXT NOT NULL,
                       password TEXT NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--DOWN
DROP TABLE users;
