package domain

import "github.com/streadway/amqp"

type AMQP interface {
	NewPublisher(exchange, queue, routingKey string) (Publisher, error)
	NewConsumer(exchange, queue, routingKey string) (Consumer, error)
}

type Publisher interface {
	Publish(message []byte) error
}

type Consumer interface {
	Consume(handler func(message amqp.Delivery))
}
