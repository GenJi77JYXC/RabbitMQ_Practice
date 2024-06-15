package main

import rabbitmq "demo/RabbitMQ"

func main() {
	exchange2 := rabbitmq.NewRabbitMQTopic("topic_exchange2", "topic.*.two")
	exchange2.RecieveTopic()
}
