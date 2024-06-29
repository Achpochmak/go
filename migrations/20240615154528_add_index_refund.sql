-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_orders_refund ON orders (refund);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_orders_refund;
-- +goose StatementEnd
