package main

import "WareSeckill/RabbitMQ"

func main() {
	rabbitmqAll := RabbitMQ.NewRabbitMQTopic("charterTopic", "charter.#")
	rabbitmqAll.ConsumerTopic()
}