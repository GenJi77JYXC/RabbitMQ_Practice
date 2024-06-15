package main

import (
	rabbitmq "demo/RabbitMQ"
	"strconv"
)

func main() {
	exchange1 := rabbitmq.NewRabbitMQTopic("topic_exchange1", "topic.hello.one")
	exchange2 := rabbitmq.NewRabbitMQTopic("topic_exchange2", "topic.hello.two")
	exchange3 := rabbitmq.NewRabbitMQTopic("topic_exchange2", "topic.nihao.two")
	exchange4 := rabbitmq.NewRabbitMQTopic("topic_exchange1", "topic.nihao.nihao")
	for i := 0; i < 10; i++ {
		msg := "hello world" + strconv.Itoa(i)
		exchange1.PublishTopic(msg + "   topic.hello.one")
		exchange2.PublishTopic(msg + "   topic.hello.two")
		exchange3.PublishTopic(msg + "   topic.nihao.two")
		exchange4.PublishTopic(msg + "   topic.nihao.nihao")
		// time.Sleep(1 * time.Second)
		// fmt.Println(i)
	}
}
