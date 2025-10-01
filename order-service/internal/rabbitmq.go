package internal

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	ch *amqp.Channel
}

func NewRabbitMQ(url string) *RabbitMQ {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("RabbitMQ not available: %v", err)
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return &RabbitMQ{ch: ch}
}

func (r *RabbitMQ) Publish(event string, payload interface{}) error {
	if r == nil {
		return nil // kalau RabbitMQ nggak jalan, skip aja
	}
	body, _ := json.Marshal(payload)
	return r.ch.Publish("", event, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
