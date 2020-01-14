package main

import "WareSeckill/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("charterPubSub")
	rabbitmq.ConsumerSub()
}
