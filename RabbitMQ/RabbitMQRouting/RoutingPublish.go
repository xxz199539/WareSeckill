package main

import (
	"WareSeckill/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmqOne := RabbitMQ.NewRabbitMQRouting("charterRouting", "routing_one")
	rabbitmqTwo:= RabbitMQ.NewRabbitMQRouting("charterRouting", "routing_two")
	for i := 0; i <= 10; i ++ {
		rabbitmqOne.PublishRouting("Hello, World！ One" + strconv.Itoa(i))
		rabbitmqTwo.PublishRouting("Hello, World！ Two" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println("发送第" +strconv.Itoa(i)+ "条消息")
	}
}
