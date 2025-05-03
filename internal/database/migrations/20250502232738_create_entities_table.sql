-- +goose Up
-- +goose StatementBegin
CREATE TABLE entities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    linkedinUrl TEXT NOT NULL,
    credlyUrl TEXT NOT NULL,
    dateCreated TIMESTAMP NOT NULL,
    dateUpdated TIMESTAMP NOT NULL,
    tags TEXT[] NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE entities;
-- +goose StatementEnd
