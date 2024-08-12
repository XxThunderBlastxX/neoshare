-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
    key TEXT NOT NULL PRIMARY KEY ,
    name TEXT NOT NULL ,
    size INT NOT NULL ,
    last_modified TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files;
-- +goose StatementEnd
