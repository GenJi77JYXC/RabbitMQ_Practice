package main

import (
	rab "demo/rabbitmq"
	"fmt"
	"strconv"
)

func main() {
	exchange1 := rab.NewRabbitMQRouting("exchange1", "exchange_1")
	exchange2 := rab.NewRabbitMQRouting("exchange2", "exchange_2")

	for i := 0; i < 10; i++ {
		msg := "Hello World---" + strconv.Itoa(i)
		if i%2 == 0 {
			exchange1.PublishRouting(msg)
		} else {
			exchange2.PublishRouting(msg)
		}
		// exchange1.PublishRouting(msg)
		// exchange2.PublishRouting(msg)
		// time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
