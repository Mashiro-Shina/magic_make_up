package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// simple mode
type RabbitMQ struct {
	Conn *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("connect rabbitmq failed: %v\n", err)
		return nil
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Printf("create channel failed: %v\n", err)
		return nil
	}
	return &RabbitMQ{
		Conn:      conn,
		Channel:   channel,
	}
}

func (r *RabbitMQ) Destroy() {
	r.Channel.Close()
	r.Conn.Close()
}

func (r *RabbitMQ) PublishComment(data []byte) {
	_, err := r.Channel.QueueDeclare(
			commentQueue,
			false,
			false,
			false,
			false,
			nil,
		)
	if err != nil {
		log.Printf("declare comment queue failed: %v\n", err)
		return
	}

	err = r.Channel.Publish(
		"",
		commentQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		log.Printf("publish comment failed: %v\n", err)
		return
	}
}

func (r *RabbitMQ) ConsumeComment() {
	_, err := r.Channel.QueueDeclare(
			commentQueue,
			false,
			false,
			false,
			false,
			nil,
		)
	if err != nil {
		log.Printf("declare comment queue failed: %v\n", err)
		return
	}

	msg, err := r.Channel.Consume(
			commentQueue,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
	if err != nil {
		log.Printf("consume comment queue failed: %v\n", err)
		return
	}

	forever := make(chan bool)
	go func() {
		for data := range msg {
			fmt.Println(data.Body)
		}
	}()

	<-forever
}
