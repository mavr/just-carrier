package rmq

import (
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const (
	reconnectDelay = 1 * time.Second
)

type Config struct {
	ConnectionString string
}

type rmq struct {
	connectionString string
	log              *zap.Logger
}

func New(cfg Config, log *zap.Logger) *rmq {
	return &rmq{
		log:              log,
		connectionString: cfg.ConnectionString,
	}
}

func connect(connectionString, exchange, queue, routingKey string) (*amqp.Channel, error) {
	connection, err := amqp.DialConfig(connectionString, amqp.Config{})
	if err != nil {
		return nil, errors.Wrapf(err, "dialing %s failed", connectionString)
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, errors.Wrapf(err, "opening channel failed")
	}

	err = channel.ExchangeDeclare(exchange, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "declare exchange %s failed", exchange)
	}

	_, err = channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "declare queue %s failed", queue)
	}

	err = channel.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "bind queue %s to exchange %s (routing key %s) failed", queue, exchange, routingKey)
	}

	return channel, nil
}
