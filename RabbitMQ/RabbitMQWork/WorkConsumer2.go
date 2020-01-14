package main

import "WareSeckill/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("charter")
	rabbitmq.ConsumerSimple()
}
