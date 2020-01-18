package repositories

import (
	"WareSeckill/models"
	"log"
)

type IOrderRepository interface {
	Insert(*models.Order)(int64, error)
	Delete(int64) bool
	Update(*models.Order) error
	SelectOneById(int64) (*models.Order, error)
	SelectAll() ([]*models.Order, error)
	SelectManyWithInfo()(map[int]map[string]string, error)
}

type OrderManagerRepository struct {
}

func NewOrderRepository() IOrderRepository {
	return &OrderManagerRepository{}
}

func(o *OrderManagerRepository)Insert(order *models.Order) (int64, error) {
	newOrder := new(models.Order)
	newOrder.OrderStatus = order.OrderStatus
	newOrder.ProductId = order.ProductId
	newOrder.UserId = order.UserId
	affect, err := engine.Insert(&newOrder)
	if err != nil {
		log.Fatalf("insert order failed: %v", err)
		return 0, nil
	}
	return affect, nil
}

func(o *OrderManagerRepository)Delete(Id int64) bool {
	order := new(models.Order)
	_, err := engine.Id(Id).Get(&order)
	if err != nil {
		log.Fatalf("insert order failed: %v", err)
		return false
	}
	_, err = engine.Id(Id).Delete(&order)
	if err != nil {
		log.Fatalf("insert order failed: %v", err)
		return false
	}
	return true
}

func(o *OrderManagerRepository)Update(order *models.Order) error{
	newOrder := new(models.Order)
	newOrder.UserId = order.UserId
	newOrder.ProductId = order.ProductId
	newOrder.OrderStatus = order.OrderStatus
	_, err := engine.Id(order.Id).Update(&newOrder)
	if err != nil {
		log.Fatalf("update order failed: %v", err)
		return err
	}
	return nil
}

func(o *OrderManagerRepository)SelectOneById(Id int64)(*models.Order, error) {
	order := new(models.Order)
	_, err := engine.Where("id = ?", Id).Get(&order)
	if err != nil {
		log.Fatalf("get order by id failed: %v", err)
		return &models.Order{}, err
	}
	return order, nil
}

func(o *OrderManagerRepository)SelectAll()([]*models.Order, error){
	orders := make([]*models.Order, 0)
	err := engine.Find(&orders)
	if err != nil {
		log.Fatalf("get all order failed: %v", err)
		return nil, err
	}
	return orders, nil
}

func(o *OrderManagerRepository)SelectManyWithInfo()(OrderMap map[int]map[string]string, err error) {
	orderArray := make(map[int]map[string]string)
	err = engine.Join("INNER", "product", "product.id = order.product_id").Find(&orderArray)
	if err != nil {
		log.Fatalf("get all product failed: %v", err)
		return nil, err
	}
	return orderArray, nil
}