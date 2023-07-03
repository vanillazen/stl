--UP
CREATE TABLE users (
                       id TEXT PRIMARY KEY,
                       username TEXT NOT NULL,
                       name TEXT NOT NULL,
                       email TEXT NOT NULL,
                       password TEXT NOT NULL
);

--DOWN
DROP TABLE users;
