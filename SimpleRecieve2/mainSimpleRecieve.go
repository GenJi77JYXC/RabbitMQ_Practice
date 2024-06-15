package main

import "demo/rabbitmq"

func main() {
	rabbitmq := rabbitmq.NewRabbitMQSimeple("" + "simple_mode")
	rabbitmq.ConsumeSimple()

}
