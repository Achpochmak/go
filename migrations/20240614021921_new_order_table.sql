-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "orders" (
    id int PRIMARY KEY,
    id_receiver int,
    delivered BOOLEAN,
    refund BOOLEAN,
    delivered_at TIMESTAMP,
    created_at TIMESTAMP,
    storage_time TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
