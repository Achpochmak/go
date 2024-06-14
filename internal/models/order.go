package models

import (
	"time"
)

type ID int

// Заказ
type Order struct {
	ID             ID        // id заказа
	ID_receiver    ID        // id получателя
	Storage_time   time.Time // время хранения
	Delivered      bool      // доставлен ли заказ
	Refund         bool      // был ли возврат
	Delivered_time time.Time // время доставки
	Created_at     time.Time // время создания
	Hash           string    // хеш
}



