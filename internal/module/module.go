package module

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
	"HOMEWORK-1/pkg/hash"
	"errors"
	"sort"
	"time"
)

type Storage interface {
	AddOrder(Order models.Order) error
	ListOrder() ([]models.Order, error)
	ReWrite(Orders []models.Order) error
	FindOrder(Id models.Id) (models.Order, error)
	UpdateOrder(Order models.Order) ( error)
}

type Deps struct {
	Storage Storage
}

type Module struct {
	Deps
}

// NewModule .. TODO сделать описание функции
func NewModule(d Deps) Module {
	return Module{Deps: d}
}

//Добавить заказ
func (m Module) AddOrder(Order models.Order) error {
	return m.Storage.AddOrder(Order)
}

//Список заказов
func (m Module) ListOrder() ([]models.Order, error) {
	return m.Storage.ListOrder()
}

//Удалить заказ
func (m Module) DeleteOrder(Order models.Order) error {
	orders, err := m.Storage.ListOrder()
	if err != nil {
		return err
	}
	
	set := make(map[models.Id]models.Order, len(orders))
	for _, order := range orders {
		set[order.Id] = order
	}
	
	_, ok := set[Order.Id]
	if !ok {
		return nil
	}

	if time.Now().Before(Order.Storage_time){
		return errors.New("время хранения не окончилось")
	}
	if Order.Delivered{
		return customErrors.ErrDelivered
	}
	delete(set, Order.Id)

	newOrders := make([]models.Order, 0, len(set))
	for _, value := range set {
		newOrders = append(newOrders, value)
	}
	return m.Storage.ReWrite(newOrders)
}

//Доставка заказа
func (m Module) DeliverOrder(order_ids []int, id_receiver int) ([]models.Order, error){

	set := []models.Order{}
	for _,id:=range order_ids{
		order, err:=m.Storage.FindOrder(models.Id(id))
		if err!=nil{
			return nil, customErrors.ErrOrderNotFound
		}
		if !time.Now().Before(order.Storage_time){
			return nil,errors.New("время хранения окончилось")
		}
		if order.Delivered{
			return nil, customErrors.ErrDelivered
		}
		if order.Id_receiver!=models.Id(id_receiver){
			return nil, customErrors.ErrWrongReceiver
		}
		
		set=append(set, order)
	}

	for _, order := range set {
		order.Delivered = true
		order.Hash = hash.GenerateHash()
		order.Delivered_time=time.Now()
		if err := m.Storage.UpdateOrder(order); err != nil {
			return nil, customErrors.ErrNotUpdated
		}
	}
	return set, nil
}

//Поиск заказов по получателю
func (m Module) OrdersByCustomer(id_receiver int, amount int) ([]models.Order, error){

	orders, err := m.Storage.ListOrder()
	if err != nil {
		return nil, err
	}
	set:=[]models.Order{}
	for _,order:=range orders{
		if order.Id_receiver==models.Id(id_receiver) {
			set=append(set, order)
		}
	}

	sort.Slice(set, func(i, j int) bool {
		return set[j].Created_at.Before(set[i].Created_at)
	})

	if amount>0{
		return set[0:amount], nil
	}
	return set, nil
}


//Поиск заказа
func (m Module) FindOrder(n models.Id) (models.Order, error) {
	return m.Storage.FindOrder(n)
}

//Возврат заказа
func (m Module) Refund(id int, id_receiver int) (error){

		order, err:=m.Storage.FindOrder(models.Id(id))
		if err!=nil{
			return customErrors.ErrOrderNotFound
		}
		if !time.Now().Before(order.Delivered_time.Add(48 * time.Hour)){
			return errors.New("время возврата истекло")
		}
		if !order.Delivered{
			return errors.New("заказ не отдали")
		}
		if order.Id_receiver!=models.Id(id_receiver){
			return customErrors.ErrWrongReceiver
		}
	
		order.Delivered = false
		order.Refund = true
		order.Hash = hash.GenerateHash()
		if err := m.Storage.UpdateOrder(order); err != nil {
			return customErrors.ErrNotUpdated
		}
	
	return nil
}

//Список возвратов
func (m Module) ListRefund() ([]models.Order, error) {
	orders, err := m.Storage.ListOrder()
	refunds:=[]models.Order{}
	for _,order:=range orders{
		if order.Refund{
			refunds=append(refunds, order)
		}
	}
	return refunds, err

}
