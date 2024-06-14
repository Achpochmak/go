package schema

import (
	"database/sql"
)

type OrderInfo struct {
	ID           int          `db:"id"`
	ID_receiver  int          `db:"id_receiver"`
	Storage_time sql.NullTime `db:"storage_time"`
	Delivered    bool         `db:"delivered"`
	Refund       bool         `db:"refund"`
	CreatedAt    sql.NullTime `db:"created_at"`
	DeliveredAt  sql.NullTime `db:"delivered_at"`
}


