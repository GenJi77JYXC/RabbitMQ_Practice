package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// 连接信息
// 这个信息是固定不变的amqp://
// 固定参数后面两个是用户名密码ip地址端口号Virtual Host 这里的simple_mode是RabbitMQ的默认虚拟主机名称，而不是指Exchange交换器
const MQURL = "amqp://genji:Gen520@127.0.0.1:5672/simple_mode"

// RabbitMQ 结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	// bind key 名称 绑定，用于消息队列和交换机之间的关联。
	Key string
	// 连接信息
	Mqurl string
}

// NewRabbitMQ 新建一个RabbitMQ实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
}

// 断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}

}

// 创建简单模式下的RabbitMQ实例
func NewRabbitMQSimeple(queueName string) *RabbitMQ {
	// 新建RabbitMQ实例
	rabbitMQ := NewRabbitMQ(queueName, "", "")
	var err error
	// 获取connection
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.Mqurl)
	rabbitMQ.failOnError(err, "Failed to connect to RabbitMQ")
	// 获取channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnError(err, "Failed to open a channel")

	return rabbitMQ
}

// simple模式下队列生产
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName, //队列名
		false,       //是否持久化
		false,       //是否排他
		false,       //是否自动删除
		false,       //是否阻塞
		nil,         //额外属性
	)
	if err != nil {
		fmt.Println("队列创建失败" + err.Error())
	}
	//调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,  //交换机名
		r.QueueName, //队列名
		false,       //mandatory 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,       //immediate 如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// simple模式下队列消费
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName, //队列名
		false,       //是否持久化
		false,       //是否排他
		false,       //是否自动删除
		false,       //是否阻塞
		nil,         //额外属性
	)
	if err != nil {
		fmt.Println("队列创建失败" + err.Error())
	}
	// 接收消息
	msgs, err := r.channel.Consume(
		q.Name, //队列名
		// consumer 用来标识消费者的身份，当多个消费者在同一个channel上消费同一个队列时，需要用不同的consumer区分
		"",    //consumer
		false, //autoAck 是否自动确认，如果为true，则表示消费者在接收到消息后，不需要再发送确认给RabbitMQ，否则需要调用channel.Ack确认消息
		false, //exclusive 是否排他(是否独有)，如果为true，则表示当前消费者只能接收到该队列的消息，不能接收到其他队列的消息
		false, //noLocal 是否不接收本机的消息
		false, //noWait 是否不等待，如果为true，则表示如果没有消息可以消费，则不会阻塞，立即返回nil
		nil,   //args 其他属性

	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启用协程处理消息
	go func() {
		// 消息处理逻辑， 可以自行设计
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// 确认消息
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 订阅模式下创建RabbitMQ实例
func NewRabbitMQPubSub(exchange string) *RabbitMQ {
	rabbitMQ := NewRabbitMQ("", exchange, "")
	var err error
	// 获取connection
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.Mqurl)
	rabbitMQ.failOnError(err, "Failed to connect to RabbitMQ")
	// 获取channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnError(err, "Failed to open a channel")

	return rabbitMQ
}

// 订阅模式下队列生产
func (r *RabbitMQ) PublishPub(message string) {
	//1. 尝试创建交换机，如果交换机不存在会自动创建，存在则跳过创建
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名
		"fanout",   //交换机类型 交换类型定义了交换的功能，即消息如何通过它路由。一旦声明了交换，就不能更改其类型。常见的类型有“direct”、“fanout”、“topic”和“headers”。
		// Fanout是持久的，因此绑定到这些预先声明的交换器的队列也必须是持久的。
		true,  //是否持久化
		false, //是否自动删除
		false, //是否是内部交换 当你想要实现内部交换时，它不应该被暴露给broker的用户 true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false, //是否阻塞 当noWait为true时，无需等待服务器的确认即可声明。通道可能由于错误而关闭。
		nil,   //额外属性
	)

	r.failOnError(err, "Failed to declare an exchange")

	// 2. 发送消息
	err = r.channel.Publish(
		r.Exchange, //交换机名
		"",         //routing key 路由键，用于指定该消息应该投递到哪个队列。如果为空，则消息会投递到默认队列。
		false,      //mandatory 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,      //immediate 如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	r.failOnError(err, "Failed to publish a message")

}

// 订阅模式下消费端代码
func (r *RabbitMQ) RecieveSub() {
	//1. 尝试创建交换机，如果交换机不存在会自动创建，存在则跳过创建
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名
		"fanout",   //交换机类型 交换类型定义了交换的功能，即消息如何通过它路由。一旦声明了交换，就不能更改其类型。常见的类型有“direct”、“fanout”、“topic”和“headers”。
		// Fanout是持久的，因此绑定到这些预先声明的交换器的队列也必须是持久的。
		true,  //是否持久化
		false, //是否自动删除
		false, //是否是内部交换 当你想要实现内部交换时，它不应该被暴露给broker的用户 true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false, //是否阻塞 当noWait为true时，无需等待服务器的确认即可声明。通道可能由于错误而关闭。
		nil,   //额外属性
	)
	r.failOnError(err, "Failed to declare an exchange")
	// 2. 尝试创建队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		"",    //队列名为空，则RabbitMQ会随机生成一个队列名
		false, //是否持久化
		false, //是否排他
		true,  //是否自动删除
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnError(err, "Failed to declare a queue")
	// 3. 绑定队列到交换机
	err = r.channel.QueueBind(
		q.Name, //队列名
		//在pub/sub模式下，这里的key要为空
		"",         //routing key 路由键，用于指定该消息应该投递到哪个队列。如果为空，则消息会投递到默认队列。
		r.Exchange, //交换机名
		false,      //noWait 如果为true，则表示如果没有消息可以消费，则不会阻塞，立即返回nil
		nil,        //额外属性
	)
	r.failOnError(err, "Failed to bind a queue")

	// 4. 消费消息

	msgs, err := r.channel.Consume(
		q.Name, //队列名
		// consumer 用来标识消费者的身份，当多个消费者在同一个channel上消费同一个队列时，需要用不同的consumer区分
		"",    //consumer
		false, //autoAck 是否自动确认，如果为true，则表示消费者在接收到消息后，不需要再发送确认给RabbitMQ，否则需要调用channel.Ack确认消息
		false, //exclusive 是否排他(是否独有)，如果为true，则表示当前消费者只能接收到该队列的消息，不能接收到其他队列的消息
		false, //noLocal 是否不接收本机的消息
		false, //noWait 是否不等待，如果为true，则表示如果没有消息可以消费，则不会阻塞，立即返回nil
		nil,   //args 其他属性
	)
	r.failOnError(err, "Failed to consume a message")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 路由模式
// 创建路由模式下的RabbitMQ实例
func NewRabbitMQRouting(exchange string, key string) *RabbitMQ {
	rabbitMQ := NewRabbitMQ("", exchange, key)
	var err error
	// 获取connection
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.Mqurl)
	rabbitMQ.failOnError(err, "Failed to connect to RabbitMQ")
	// 获取channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnError(err, "Failed to open a channel")

	return rabbitMQ
}

// 路由模式下队列生产
func (r *RabbitMQ) PublishRouting(message string) {
	//1. 尝试创建交换机，如果交换机不存在会自动创建，存在则跳过创建
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名
		"direct",   //交换机类型 交换类型定义了交换的功能，即消息如何通过它路由。一旦声明了交换，就不能更改其类型。常见的类型有“direct”、“fanout”、“topic”和“headers”。
		// Direct是持久的，因此绑定到这些预先声明的交换器的队列也必须是持久的。
		true,  //是否持久化
		false, //是否自动删除
		false, //是否是内部交换 当你想要实现内部交换时，它不应该被暴露给broker的用户 true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false, //是否阻塞 当noWait为true时，无需等待服务器的确认即可声明。通道可能由于错误而关闭。
		nil,   //额外属性
	)
	r.failOnError(err, "Failed to declare an exchange")
	// 2. 发送消息
	err = r.channel.Publish(
		r.Exchange, //交换机名
		r.Key,      //routing key 路由键，用于指定该消息应该投递到哪个队列。如果为空，则消息会投递到默认队列。
		false,      //mandatory 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,      //immediate 如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	r.failOnError(err, "Failed to publish a message")
}

// 路由模式下消费端代码
func (r *RabbitMQ) RecieveRouting() {
	//1. 尝试创建交换机，如果交换机不存在会自动创建，存在则跳过创建
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名
		"direct",   //交换机类型 交换类型定义了交换的功能，即消息如何通过它路由。一旦声明了交换，就不能更改其类型。常见的类型有“direct”、“fanout”、“topic”和“headers”。
		// Direct是持久的，因此绑定到这些预先声明的交换器的队列也必须是持久的。
		true,  //是否持久化
		false, //是否自动删除
		false, //是否是内部交换 当你想要实现内部交换时，它不应该被暴露给broker的用户 true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false, //是否阻塞 当noWait为true时，无需等待服务器的确认即可声明。通道可能由于错误而关闭。
		nil,   //额外属性
	)
	r.failOnError(err, "Failed to declare an exchange")
	// 2. 尝试创建队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		"",    //队列名为空，则RabbitMQ会随机生成一个队列名
		false, //是否持久化
		false, //是否排他
		true,  //是否自动删除
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnError(err, "Failed to declare a queue")
	// 3. 绑定队列到交换机
	err = r.channel.QueueBind(
		q.Name,     //队列名
		r.Key,      //routing key 路由键，用于指定该消息应该投递到哪个队列。如果为空，则消息会投递到默认队列。
		r.Exchange, //交换机名
		false,      //noWait 如果为true，则表示如果没有消息可以消费，则不会阻塞，立即返回nil
		nil,        //额外属性
	)
	r.failOnError(err, "Failed to bind a queue")

	// 4. 消费消息

	msgs, err := r.channel.Consume(
		q.Name, //队列名
		// consumer 用来标识消费者的身份，当多个消费者在同一个channel上消费同一个队列时，需要用不同的consumer区分
		"",    //consumer
		false, //autoAck 是否自动确认，如果为true，则表示消费者在接收到消息后，不需要再发送确认给RabbitMQ，否则需要调用channel.Ack确认消息
		false, //exclusive 是否排他(是否独有)，如果为true，则表示当前消费者只能接收到该队列的消息，不能接收到其他队列的消息
		false, //noLocal 是否不接收本机的消息
		false, //noWait 是否不等待，如果为true，则表示如果没有消息可以消费，则不会阻塞，立即返回nil
		nil,   //args 其他属性
	)
	r.failOnError(err, "Failed to consume a message")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 话题模式
// 创建话题模式下的RabbitMQ实例
func NewRabbitMQTopic(exchange string, routingKey string) *RabbitMQ {
	rabbitMQ := NewRabbitMQ("", exchange, routingKey)
	var err error
	// 获取connection
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.Mqurl)
	rabbitMQ.failOnError(err, "Failed to connect to RabbitMQ")
	// 获取channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnError(err, "Failed to open a channel")

	return rabbitMQ

}

// 话题模式下队列生产

func (r *RabbitMQ) PublishTopic(message string) {
	//1. 尝试创建交换机，如果交换机不存在会自动创建，存在则跳过创建
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名
		"topic",    //交换机类型 交换类型定义了交换的功能，即消息如何通过它路由。一旦声明了交换，就不能更改其类型。常见的类型有“direct”、“fanout”、“topic”和“headers”。
		// Topic是持久的，因此绑定到这些预先声明的交换器的队列也必须是持久的。
		true,  //是否持久化
		false, //是否自动删除
		false, //是否是内部交换 当你想要实现内部交换时，它不应该被暴露给broker的用户 true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false, //是否阻塞 当noWait为true时，无需等待服务器的确认即可声明。通道可能由于错误而关闭。
		nil,   //额外属性
	)
	r.failOnError(err, "Failed to declare an exchange")
	// 2. 发送消息
	err = r.channel.Publish(
		r.Exchange, //交换机名
		r.Key,      //routing key 路由键，用于指定该消息应该投递到哪个队列。如果为空，则消息会投递到默认队列。
		false,      //mandatory 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,      //immediate 如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	r.failOnError(err, "Failed to publish a message")
}

// 话题模式下消费端代码
// 要注意key,规则
// 其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
// 匹配 genji.* 表示匹配 genji.hello
// genji.hello.one需要用kuteng.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
	//1. 尝试创建交换机，如果交换机不存在会自动创建，存在则跳过创建
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名
		"topic",    //交换机类型 交换类型定义了交换的功能，即消息如何通过它路由。一旦声明了交换，就不能更改其类型。常见的类型有“direct”、“fanout”、“topic”和“headers”。
		// Topic是持久的，因此绑定到这些预先声明的交换器的队列也必须是持久的。
		true,  //是否持久化
		false, //是否自动删除
		false, //是否是内部交换 当你想要实现内部交换时，它不应该被暴露给broker的用户 true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false, //是否阻塞 当noWait为true时，无需等待服务器的确认即可声明。通道可能由于错误而关闭。
		nil,   //额外属性
	)
	r.failOnError(err, "Failed to declare an exchange")
	// 2. 尝试创建队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		"",    //队列名为空，则RabbitMQ会随机生成一个队列名
		false, //是否持久化
		false, //是否排他
		true,  //是否自动删除
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnError(err, "Failed to declare a queue")
	// 3. 绑定队列到交换机
	err = r.channel.QueueBind(
		q.Name,     //队列名
		r.Key,      //routing key 路由键，用于指定该消息应该投递到哪个队列。如果为空，则消息会投递到默认队列。
		r.Exchange, //交换机名
		false,      //noWait 如果为true，则表示如果没有消息可以消费，则不会阻塞，立即返回nil
		nil,        //额外属性
	)
	r.failOnError(err, "Failed to bind a queue")

	// 4. 消费消息

	msgs, err := r.channel.Consume(
		q.Name, //队列名
		// consumer 用来标识消费者的身份，当多个消费者在同一个channel上消费同一个队列时，需要用不同的consumer区分
		"",    //consumer
		false, //autoAck 是否自动确认，如果为true，则表示消费者在接收到消息后，不需要再发送确认给RabbitMQ，否则需要调用channel.Ack确认消息
		false, //exclusive 是否排他(是否独有)，如果为true，则表示当前消费者只能接收到该队列的消息，不能接收到其他队列的消息
		false, //noLocal 是否不接收本机的消息
		false, //noWait 是否不等待，如果为true，则表示如果没有消息可以消费，则不会阻塞，立即返回nil
		nil,   //args 其他属性
	)
	r.failOnError(err, "Failed to consume a message")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
