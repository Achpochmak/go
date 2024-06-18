package models

import (
	"time"
)

type ID int

// Заказ
type Order struct {
	ID           ID        // id заказа
	IDReceiver   ID        // id получателя
	StorageTime  time.Time // время хранения
	Delivered    bool      // доставлен ли заказ
	Refund       bool      // был ли возврат
	DeliveryTime time.Time // время доставки
	CreatedAt    time.Time // время создания
	WeightKg     float64   //вес
	Price        float64   //цена
	Packaging    Packaging //упаковка
}
