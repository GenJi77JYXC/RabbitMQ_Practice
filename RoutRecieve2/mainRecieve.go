package main

import "demo/rabbitmq"

func main() {
	exchange2 := rabbitmq.NewRabbitMQRouting("exchange2", "exchange_2")
	exchange2.RecieveRouting()
}
