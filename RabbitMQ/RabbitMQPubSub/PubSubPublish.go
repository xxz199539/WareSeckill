package main

import (
	"WareSeckill/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main()  {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("charterPubSub")
	for i := 0; i <= 100; i++ {
		rabbitmq.PublishPub("Hello World!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println("生产者生成第" +strconv.Itoa(i)+"条消息")
	}
}
