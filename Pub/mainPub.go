package main

import (
	rabbitmq "demo/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQPubSub("" + "newProduct")
	for i := 0; i < 5; i++ {
		rabbitmq.PublishPub("订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		fmt.Println("订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		time.Sleep(1 * time.Second)
	}
}
