-- +goose Up
CREATE SCHEMA IF NOT EXISTS template;

-- +goose Down
DROP SCHEMA IF EXISTS template;