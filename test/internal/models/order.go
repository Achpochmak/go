package models
import (
	"time"
)

type Id int

// Заказ
type Order struct {
	Id      Id   						// id заказа
	Id_receiver Id						// id получателя
	Storage_time time.Time 				// время хранения
	Delivered bool 						// доставлен ли заказ
	Refund bool 						// был ли возврат
	Delivered_time time.Time 			// время доставки
	Created_at time.Time 				// время создания
	Hash string							// хеш
}

