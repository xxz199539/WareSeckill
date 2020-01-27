package main

import (
	"WareSeckill/RabbitMQ"
	"WareSeckill/repositories"
	"WareSeckill/services"
)

func main() {
	productManager := repositories.NewProductManager()
	productService := services.NewProductService(productManager)
	orderManager := repositories.NewOrderRepository()
	orderService := services.NewOrderService(orderManager)
	rabbitMq := RabbitMQ.NewRabbitMQSimple("BoomShakalaka")
	rabbitMq.ConsumerSimpleByService(orderService, productService)
	defer rabbitMq.Destroy()
}
