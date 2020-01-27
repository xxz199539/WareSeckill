package repositories

import (
	"WareSeckill/common"
	"WareSeckill/models"
	"log"
	"time"
)

type IOrderRepository interface {
	Insert(*models.Order)(int64, error)
	Delete(int64) bool
	Update(*models.Order) error
	SelectOneById(int64) (*models.Order, error)
	SelectAll() ([]*models.Order, error)
	SelectManyWithInfo()([]*models.OrderGroup, error)
}

type OrderManagerRepository struct {
}

func NewOrderRepository() IOrderRepository {
	return &OrderManagerRepository{}
}

func(o *OrderManagerRepository)Insert(order *models.Order) (int64, error) {
	newOrder := &models.Order{
		UserId:      order.UserId,
		ProductId:   order.ProductId,
		OrderStatus: order.OrderStatus,
		CreateTime:  time.Now(),
	}
	affect, err := common.Engine.Insert(newOrder)
	if err != nil {
		log.Fatalf("insert order failed: %v", err)
		return 0, nil
	}
	return affect, nil
}

func(o *OrderManagerRepository)Delete(Id int64) bool {
	order := new(models.Order)
	_, err := common.Engine.Id(Id).Get(order)
	if err != nil {
		log.Fatalf("insert order failed: %v", err)
		return false
	}
	_, err = common.Engine.Id(Id).Delete(&order)
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
	_, err := common.Engine.Id(order.Id).Update(&newOrder)
	if err != nil {
		log.Fatalf("update order failed: %v", err)
		return err
	}
	return nil
}

func(o *OrderManagerRepository)SelectOneById(Id int64)(*models.Order, error) {
	order := new(models.Order)
	_, err := common.Engine.Where("id = ?", Id).Get(order)
	if err != nil {
		log.Fatalf("get order by id failed: %v", err)
		return &models.Order{}, err
	}
	return order, nil
}

func(o *OrderManagerRepository)SelectAll()([]*models.Order, error){
	orders := make([]*models.Order, 0)
	err := common.Engine.Find(orders)
	if err != nil {
		log.Fatalf("get all order failed: %v", err)
		return nil, err
	}
	return orders, nil
}

func(o *OrderManagerRepository)SelectManyWithInfo()(OrderMap []*models.OrderGroup, err error) {
	orderArray := make([]*models.OrderGroup, 0)
	err = common.Engine.Join("INNER", "product", "product.id = order.product_id").Find(orderArray)
	//err = common.Engine.SQL("Select o.id,p.product_name,o.order_status From ware_seckill.order as o left join product as p on o.product_id=p.id").Find(&orderArray)
	if err != nil {
		log.Fatalf("get all order failed: %v", err)
		return nil, err
	}
	return orderArray, nil
}