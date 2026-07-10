-- +goose Up
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(64) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE roles;
