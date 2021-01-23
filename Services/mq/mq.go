package mq

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
)

type CallBack func(msg string)

func Connect() (*amqp.Connection, error) {
	return amqp.Dial("amqp://guest:guest@127.0.0.1:15672")
}

// 发送端函数
//exchange 交换机名称  queueName队列名称 body内容
func Publish(exchange, queueName, body string) error {
	//建立连接
	conn, err := Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	//创建一个通道 channel
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	//创建队列
	q, err := channel.QueueDeclare(
		queueName,
		false, //是否持久化
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
		ContentType: "text/plain", //发布内容类型
		Body:        []byte(body),
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
		false, //是否持久化
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	//接收消息
	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			s := BytesToString(&(d.Body))
			callback(*s)
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
