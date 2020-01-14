package main

import (
	"WareSeckill/RabbitMQ"
"fmt"
)

func main(){
	rabbitmq := RabbitMQ.NewRabbitMQSimple("Simple")
	rabbitmq.PublishSimple("Hello,World!")
	fmt.Println("发送成功!")
}
