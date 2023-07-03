--UP
CREATE TABLE tags (
                       id TEXT PRIMARY KEY,
                       name TEXT NOT NULL,
                       description TEXT NOT NULL
);

--DOWN
DROP TABLE tags;
