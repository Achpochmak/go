-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD COLUMN weight_kg FLOAT,
ADD COLUMN price FLOAT,
ADD COLUMN packaging VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
DROP COLUMN weight_kg,
DROP COLUMN price,
DROP COLUMN packaging;
-- +goose StatementEnd
