-- +goose Up
-- +goose StatementBegin
ALTER TABLE entities RENAME TO contacts;
ALTER TABLE contacts ADD COLUMN contact VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE contacts DROP COLUMN contact;
ALTER TABLE contacts RENAME TO entities;
-- +goose StatementEnd
