-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "orders" (
    id int PRIMARY KEY,
    idReceiver int,
    delivered BOOLEAN,
    refund BOOLEAN,
    deliveredAt TIMESTAMP,
    createdAt TIMESTAMP,
    storageTime TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
