package main

import "demo/rabbitmq"

func main() {
	exchange3 := rabbitmq.NewRabbitMQRouting("exchange2", "exchange_2")
	exchange3.RecieveRouting() // 当调用RecieveRouting()时，如果绑定的key没有对应的队列，则会自动创建队列。如果有队列，则会将消息推送到队列中。
}
