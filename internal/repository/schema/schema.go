package schema

import (
	"database/sql"
)

type OrderInfo struct {
	ID          int          `db:"id"`
	IDReceiver  int          `db:"idreceiver"`
	StorageTime sql.NullTime `db:"storagetime"`
	Delivered   bool         `db:"delivered"`
	Refund      bool         `db:"refund"`
	CreatedAt   sql.NullTime `db:"createdat"`
	DeliveredAt sql.NullTime `db:"deliveredat"`
}
