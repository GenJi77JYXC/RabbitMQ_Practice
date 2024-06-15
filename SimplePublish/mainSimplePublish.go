package main

import (
	"demo/rabbitmq"
	"fmt"
	"time"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQSimeple("" + "simple_mode")
	// rabbitmq.PublishSimple("Hello World!")
	for i := 0; i < 10; i++ {
		go publishMsg(i, rabbitmq)
	}
	time.Sleep(3 * time.Second)
	fmt.Println("Published message to RabbitMQ")
}

func publishMsg(i int, rabbitmq *rabbitmq.RabbitMQ) {
	rabbitmq.PublishSimple(fmt.Sprintf("Hello World! %d", i))
}
