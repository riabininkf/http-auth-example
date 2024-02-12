-- +goose Up
CREATE TABLE IF NOT EXISTS template.products
(
    id         BIGSERIAL                  NOT NULL,
    name       TEXT                    NOT NULL,
    comment    TEXT,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS products_name_uidx ON template.products (name);

-- +goose Down
DROP TABLE IF EXISTS template.products;
