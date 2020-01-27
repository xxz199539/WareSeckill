package main

import (
	"WareSeckill/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("")
	rabbitmq.ConsumerSimple()
}
