package customErrors

import (
	"errors"
)

var (
	ErrIDNotFound          = errors.New("ID заказа не найден")
	ErrReceiverNotFound    = errors.New("получатель не найден")
	ErrNotUpdated          = errors.New("не удалось обновить заказ")
	ErrDelivered           = errors.New("заказ отдали")
	ErrWrongReceiver       = errors.New("не тот получатель")
	ErrOrderNotFound       = errors.New("заказ не найден")
	ErrStorageTimeEnded    = errors.New("время хранения окончилось")
	ErrStorageTimeNotEnded = errors.New("время хранения не окончилось")
	ErrNotDelivered        = errors.New("заказ не отдали")
	ErrRefundTimeEnded     = errors.New("время возврата истекло")
	ErrWorkersLessThanOne  = errors.New("количество должно быть больше 1")
	ErrWrongTimeFormat     = errors.New("неправильный формат времени")
	ErrOrderAlreadyExists  = errors.New("заказ дублируется")
)
