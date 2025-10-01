package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	url      string
	conn     *amqp.Connection
	ch       *amqp.Channel
	exchange string
}

func NewRabbitMQ(url string) *RabbitMQ {
	return &RabbitMQ{url: url, exchange: "events"}
}

func (r *RabbitMQ) Connect() error {
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	if err := ch.ExchangeDeclare(r.exchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	r.conn = conn
	r.ch = ch
	return nil
}

func (r *RabbitMQ) Publish(routingKey string, payload any) error {
	bts, _ := json.Marshal(payload)
	return r.ch.Publish(r.exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bts,
	})
}

func (r *RabbitMQ) Subscribe(queueName, bindingKey string, onMessage func(map[string]interface{})) {
	q, err := r.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := r.ch.QueueBind(q.Name, bindingKey, r.exchange, false, nil); err != nil {
		log.Fatal(err)
	}
	msgs, err := r.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for d := range msgs {
			var payload map[string]interface{}
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Println("invalid json in message", err)
				d.Nack(false, false)
				continue
			}
			onMessage(payload)
			d.Ack(false)
		}
	}()
}
