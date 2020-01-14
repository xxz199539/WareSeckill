package main

import "WareSeckill/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQRouting("charterRouting", "routing_two")
	rabbitmq.ConsumerRouting()
}