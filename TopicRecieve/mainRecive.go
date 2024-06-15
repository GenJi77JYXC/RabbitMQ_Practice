package main

import rabbitmq "demo/RabbitMQ"

func main() {
	exchange1 := rabbitmq.NewRabbitMQTopic("topic_exchange1", "#")
	exchange1.RecieveTopic()
}
