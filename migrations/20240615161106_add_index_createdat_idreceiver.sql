-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_orders_id_receiver_created_at ON orders (id_receiver, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_orders_id_receiver_created_at;
-- +goose StatementEnd