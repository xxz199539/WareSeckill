package main

import (
	"WareSeckill/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmqOne := RabbitMQ.NewRabbitMQTopic("charterTopic", "charter.topic.one")
	rabbitmqTwo := RabbitMQ.NewRabbitMQTopic("charterTopic", "charter.topic.two")

	for i := 0; i <= 10; i ++ {
		rabbitmqOne.PublishTopic("Hello, World！ One" + strconv.Itoa(i))
		rabbitmqTwo.PublishTopic("Hello, World！ Two" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println("发送第" + strconv.Itoa(i) + "条消息")
	}
}