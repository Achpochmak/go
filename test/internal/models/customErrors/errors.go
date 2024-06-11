package customErrors

import (
	"errors"
)

var (
	ErrIdNotFound = errors.New("ID заказа не найден")
	ErrReceiverNotFound = errors.New("получатель не найден")
	ErrNotUpdated = errors.New("не удалось обновить заказ")
	ErrDelivered = errors.New("заказ отдали")
	ErrWrongReceiver = errors.New("не тот получатель")
	ErrOrderNotFound = errors.New("заказ не найден")
)