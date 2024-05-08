-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    date TIMESTAMP NOT NULL,
    description TEXT,
    price NUMERIC(10, 2)
);
CREATE TABLE clubs (
   id SERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL UNIQUE,
   description TEXT,
   is_active VARCHAR(6) NOT NULL DEFAULT 'true',
   contacts TEXT,
   price NUMERIC(10, 2),
   spots_number INT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS clubs;
-- +goose StatementEnd
