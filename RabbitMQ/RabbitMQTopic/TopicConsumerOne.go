package main

import "WareSeckill/RabbitMQ"

func main() {
	rabbitmqAll := RabbitMQ.NewRabbitMQTopic("charterTopic", "charter.*.two")
	rabbitmqAll.ConsumerTopic()
}