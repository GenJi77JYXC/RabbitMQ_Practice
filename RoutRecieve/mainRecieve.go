package main

import "demo/rabbitmq"

func main() {
	exchange1 := rabbitmq.NewRabbitMQRouting("exchange1", "exchange_1")
	exchange1.RecieveRouting()
}
