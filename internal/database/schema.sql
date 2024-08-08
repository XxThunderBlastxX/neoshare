CREATE TABLE file (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL ,
    key TEXT NOT NULL ,
    size INT NOT NULL ,
    last_modified TIMESTAMP NOT NULL
);

CREATE INDEX idx_file_name ON file(name);

CREATE INDEX idx_file_key ON file(key);

CREATE INDEX idx_file_size ON file(size);

CREATE INDEX idx_last_modified ON file(last_modified);