package services

import (
	"WareSeckill/models"
	"WareSeckill/repositories"
)

type IOrderService interface {
	Insert(*models.Order)(int64, error)
	Delete(int64) bool
	Update(*models.Order) error
	SelectOneById(int64) (*models.Order, error)
	SelectAll() ([]*models.Order, error)
	SelectManyWithInfo()(map[int]map[string]string, error)
}

type OrderService struct {
	orderRepository repositories.IOrderRepository
}

func NewOrderService(repositories repositories.IOrderRepository) IOrderService {
	return &OrderService{orderRepository: repositories}
}

func (o *OrderService) Insert(order *models.Order) (int64, error) {
	return o.orderRepository.Insert(order)
}

func (o *OrderService) Delete(Id int64) bool {
	return o.orderRepository.Delete(Id)
}

func (o *OrderService) Update(order *models.Order) error {
    return o.orderRepository.Update(order)
}

func (o *OrderService)SelectOneById(Id int64)(*models.Order, error) {
	return o.orderRepository.SelectOneById(Id)
}

func (o *OrderService)SelectAll()([]*models.Order, error) {
	return o.orderRepository.SelectAll()

}

func (o *OrderService)SelectManyWithInfo()(map[int]map[string]string, error) {
	return o.orderRepository.SelectManyWithInfo()
}
