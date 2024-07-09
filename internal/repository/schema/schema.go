package schema

import (
	"database/sql"
)

type OrderInfo struct {
	ID          int          `db:"id"`
	IDReceiver  int          `db:"id_receiver"`
	StorageTime sql.NullTime `db:"storage_time"`
	Delivered   bool         `db:"delivered"`
	Refund      bool         `db:"refund"`
	CreatedAt   sql.NullTime `db:"created_at"`
	DeliveredAt sql.NullTime `db:"delivered_at"`
	WeightKg    float64      `db:"weight_kg"`
	Price       float64      `db:"price"`
	Packaging   string       `db:"packaging"`
}
