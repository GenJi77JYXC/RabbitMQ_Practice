package main

import "demo/rabbitmq"

func main() {
	rabbitmq := rabbitmq.NewRabbitMQPubSub("" + "newProduct")
	rabbitmq.RecieveSub()
}
