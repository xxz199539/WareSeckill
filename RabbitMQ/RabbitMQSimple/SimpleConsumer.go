package main

import "WareSeckill/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("Simple")
	rabbitmq.ConsumerSimple()
}
