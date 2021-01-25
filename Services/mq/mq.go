package mq

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
)

type CallBack func(msg string)

func Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	return conn, err
}

// 发送端函数
//exchange 交换机名称  queueName队列名称 body内容
func Publish(exchange, queueName, body string) error {
	//建立连接
	conn, err := Connect()
	if err != nil {
		fmt.Println("建立连接失败:", err)
		return err
	}
	defer conn.Close()

	//创建一个通道 channel
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println("创建一个通道失败:", err)
		return err
	}
	defer channel.Close()

	//创建队列
	q, err := channel.QueueDeclare(
		queueName,
		true, //是否持久化
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil
	}
	//发送消息
	err = channel.Publish(exchange, q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain", //发布内容类型
		Body:         []byte(body),
	})
	return err
}

//接收者的方法
func Consumer(exchange, queueName string, callback CallBack) {
	//建立连接
	conn, err := Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//创建一个通道 channel
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer channel.Close()

	//创建队列
	q, err := channel.QueueDeclare(
		queueName,
		true, //是否持久化
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	//接收消息 第三参数为true 表示自动应答
	msgs, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			s := BytesToString(&(d.Body))
			callback(*s)
			d.Ack(false) //手动应答
		}
	}()
	fmt.Println("waiting for messages")
	<-forever
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

// 订阅模式
// exchange 交换机
// types 类型决定了交换机是订阅模式还是路由模式还是主题模式
// 路由的key
// 内容
func PublishEx(exchange, types, routingKey, body string) error {
	//建立连接
	conn, err := Connect()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	//创建一个通道 channel
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer channel.Close()

	//创建交换机
	err = channel.ExchangeDeclare(
		exchange,
		types,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	err = channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	return err
}

//消费端
func ConsumeEx(exchange, types, routingKey string, callback CallBack) {
	//建立连接
	conn, err := Connect()
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建通道
	channel, err := conn.Channel()
	defer channel.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建交换机
	err = channel.ExchangeDeclare(exchange, types, true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建队列  临时队列
	q, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//绑定
	err = channel.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//接收信息
	msgs, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	forever := make(chan bool)
	go func() {
		for {
			for d := range msgs {
				s := BytesToString(&(d.Body))
				callback(*s)
				d.Ack(false)
			}
		}
	}()
	fmt.Printf("wait for message:\n")
	<-forever
}

//死性队列消费端
func ConsumerDlx(exchangeA, queueAName, exchangeB, queueBName string, ttl int, callback CallBack) {
	//建立连接
	conn, err := Connect()
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	//创建channel
	channel, err := conn.Channel()
	defer channel.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	//创建A交换机
	err = channel.ExchangeDeclare(
		exchangeA,
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//创建A队列
	queueA, err := channel.QueueDeclare(exchangeA, true, false, false, false, amqp.Table{
		"x-message-ttl":          ttl,
		"x-dead-letter-exchange": exchangeB,
		//"x-dead-letter-queue":"",//绑定哪个队列
		//"x-dead-letter-routing-key":"",//绑定路由的关键字
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	//A交换机和A队列绑定
	err = channel.QueueBind(queueA.Name, "", exchangeA, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建B交换机
	err = channel.ExchangeDeclare(exchangeB, "fanout", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//创建B队列
	queueB, err := channel.QueueDeclare(queueBName, true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//绑定B交换机和B队列
	err = channel.QueueBind(queueB.Name, "", exchangeB, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	msgs, err := channel.Consume(queueB.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			s := BytesToString(&(d.Body))
			callback(*s)
			d.Ack(false)
		}
	}()
	fmt.Println("waiting for message")
	<-forever
}

//死性队列生产端
func PublishDlx(exchangeA, body string) error {
	//建立连接
	conn, err := Connect()
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	//创建channel
	channel, err := conn.Channel()
	defer channel.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = channel.Publish(exchangeA, "", false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(body),
	})
	return err
}
