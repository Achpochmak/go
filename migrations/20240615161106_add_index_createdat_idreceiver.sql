-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_orders_idReceiver_createdAt ON orders (idReceiver, createdAt DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_orders_idReceiver_createdAt;
-- +goose StatementEnd