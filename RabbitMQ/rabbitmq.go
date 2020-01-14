package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const MQURL = "amqp://charter:xxz199439@127.0.0.1:5672/charter"

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	// queue name
	QueueName string
	Exchange string
	Key string
	Mqurl string
}
// 创建实例.此方法公用
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "get rabbitmq conn fail")
	rabbitmq.channel,err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "get rabbitmq channel fail")
	return rabbitmq
}


func (r *RabbitMQ) Destroy() {
	err := r.channel.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("rabbitmq close channel fail:%v", err.Error()))
	}
	err = r.conn.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("rabbitmq close conn fail:%v", err.Error()))
	}
}

// deal error function
func (r *RabbitMQ) failOnError (err error, message string){
	if err != nil{
		log.Fatal(fmt.Sprintf("%s:%s", message, err))
	}
}
// 简单模式Step1：创建简单模式下的RabbitMQ实例
func NewRabbitMQSimple (queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}
//简单模式Step2:生产者
func (r *RabbitMQ) PublishSimple(message string) {
	// 1.先申请队列，如果找不到队列则新建队列，保证消息能发送到队列中
    _, err := r.channel.QueueDeclare(
    	r.QueueName,
    	false, // 消息是否持久化
    	false, // 消息是否自动删除
    	false, // 是否具有排他性，只有自己能访问
    	false,// 是否阻塞
    	nil, // 额外属性
    	)
    if err != nil {
    	fmt.Println(err)
	}
	// 2.发送消息到队列中
	err = r.channel.Publish(
		r.Exchange, // 使用默认的交换机类型
		r.QueueName,
		false, // 如果为True，会根据exchange类型和routekey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false, // true:当exchange发送消息到队列后发现队列上没有绑定消费者，会把消息返回给生产者
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(message),
		})
	if err != nil {
		fmt.Println(err)
	}
}

func (r *RabbitMQ) ConsumerSimple() {
	// 1.先申请队列，如果找不到队列则新建队列，保证消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 消息是否持久化
		false, // 消息是否自动删除
		false, // 是否具有排他性，只有自己能访问
		false,// 是否阻塞
		nil, // 额外属性
	)
	if err != nil {
		fmt.Println(err)
	}
	// 接受消息
	msms, err := r.channel.Consume(
		r.QueueName,
		"", // 用来区分多个消费者
		true, // 是否自动应答
		false, // 是否具有排他性
		false, // true:表示不能讲同一个connection中发送的消息传递给这个connection中的消费者
		false, // true：表示阻塞
		nil, // 其他属性
		)
	if err != nil {
		log.Fatal(err)
	}
	forever := make(chan bool)
	go func() {
		for msg := range msms {
			// 实现要处理的逻辑函数
			log.Printf("Receive a messagee %s", msg.Body)
		}
		forever <- true
	}()
	log.Println("[*] waiting for message, to exit press CTRL+C")
    <- forever
}

// 订阅模式创建实例
func NewRabbitMQPubSub(exchange string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchange, "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "get rabbitmq conn fail")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "get rabbitmq channel fail")
	return rabbitmq
}

// 订阅模式生产者
func (r *RabbitMQ) PublishPub(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",// 广播类型
		true,
		false,
		false,// true:表示exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定。
		false,
		nil,
		)
	r.failOnError(err, "failed to declare an exchange")
    err = r.channel.Publish(
    	r.Exchange,
    	"",
    	false,
    	false,
    	amqp.Publishing{
    		ContentType:"text/plain",
    		Body:[]byte(message)})
    if err!= nil {
    	fmt.Println(err)
	}
}

func (r *RabbitMQ) ConsumerSub() {
	// 1.试着创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
		)
	r.failOnError(err, "failed to declare an exchange")
	// 2.试着创建队列
	q, err := r.channel.QueueDeclare(
		"", // 随机生成队列名称
		false,
		false,
		true,
		false,
		nil,
		)
	r.failOnError(err, "failed to declare an queue")
	// 3.绑定队列到交换机
	err = r.channel.QueueBind(
		q.Name,
		"", // 在pub/sub模式下，这里的key为空
		r.Exchange,
		false,
		nil,
		)
	r.failOnError(err, "failed to binded queue to exchange")
	// 4.消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
		)
	if err != nil {
		log.Println(err)
	}
	forever := make(chan bool)
	go func() {
		for msg := range messages {
			// 实现要处理的逻辑函数
			log.Printf("Receive a messagee %s", msg.Body)
		}
		forever <- true
	}()
	log.Println("[*] waiting for message, to exit press CTRL+C")
	<- forever
}

// 创建路由模式实例
func NewRabbitMQRouting(exchangeName, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "get rabbitmq conn fail")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "get rabbitmq channel fail")
	return rabbitmq
}

// 路由模式生产者
func(r *RabbitMQ) PublishRouting(message string) {
	// 1.试着创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
		)
	r.failOnError(err, "failed to declare an exchange")

	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType:"text/plain",
			Body:[]byte(message),
		})
}

// 路由模式消费者
func (r *RabbitMQ) ConsumerRouting(){
	// 1.试着创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
		)
	r.failOnError(err, "failed to declare an exchange")
	// 2.试着创建队列
	q, err := r.channel.QueueDeclare(
		"", // 随机队列名称
		false,
		false,
		true,
		false,
		nil,
		)
	r.failOnError(err, "failed to declare an queue")
    // 3.绑定队列与交换机
    err = r.channel.QueueBind(
    	q.Name,
    	r.Key,
    	r.Exchange,
    	false,
    	nil,
    	)
	r.failOnError(err, "failed to bind queue to exchange")

    messages, err := r.channel.Consume(
    	q.Name,
    	"",
    	true,
    	false,
    	false,
    	false,
    	nil,
    	)
    if err != nil {
    	log.Println(err)
	}
    // 消费消息
    forever := make(chan bool)
    go func() {
    	for msg := range messages{
			log.Printf("Receive a messagee %s", msg.Body)
		}
		forever <- true
	}()
    <- forever
}

// 创建话题模式实例
func NewRabbitMQTopic(exchangeName, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "get rabbitmq conn fail")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "get rabbitmq channel fail")
	return rabbitmq
}

// 话题模式生产者

func (r *RabbitMQ)PublishTopic(message string){
	// 1.试着创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
		)
	r.failOnError(err, "failed to declare exchange")
	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType:"text/plain",
			Body:[]byte(message),
		})
}

// 话题模式消费者。key为"*"匹配一个单词，"#"匹配多个单词（可以使零个）
// 匹配charter.*表示匹配charter.hello,但是charter.hello.one需要用charter.#才能匹配到
func (r *RabbitMQ)ConsumerTopic(){
	// 1.创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err, "failed to declare exchange")
	// 2.创建队列
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
		)
	r.failOnError(err, "failed to declare queue")
	// 3.绑定队列到交换机
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
		)
	r.failOnError(err, "failed to bind queue to exchange")
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	// 消费消息
	forever := make(chan bool)
	go func() {
		for msg := range messages{
			log.Printf("Receive a messagee %s", msg.Body)
		}
		forever <- true
	}()
	<- forever
}