package rmq

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type consumer struct {
	log *zap.Logger

	connectionString string
	queue            string
	routingKey       string
	exchange         string
}

func (c *rmq) NewConsumer(exchange, queue, routingKey string) (*consumer, error) {
	consumer := &consumer{
		log:              c.log.Named(fmt.Sprintf("consumer(%s::%s)", queue, routingKey)),
		connectionString: c.connectionString,
		queue:            queue,
		routingKey:       routingKey,
		exchange:         exchange,
	}

	return consumer, nil
}

func (c *consumer) Consume(handler func(message amqp.Delivery)) {
	for {
		messages, err := c.consume()
		if err != nil {
			c.log.Error("consume", zap.Error(err), zap.Int("delay_sec", int(reconnectDelay.Seconds())))
			time.Sleep(reconnectDelay)
			continue
		}

		c.log.Debug("starting consuming",
			zap.String("exchange", c.exchange),
			zap.String("queue", c.queue),
			zap.String("routing_key", c.routingKey),
		)


		for delivery := range messages {
			handler(delivery)
		}

		c.log.Error("reconnect", zap.Int("delay_sec", int(reconnectDelay.Seconds())))
		time.Sleep(reconnectDelay)
	}
}

func (c *consumer) consume() (<-chan amqp.Delivery, error) {
	ch, err := connect(c.connectionString, c.exchange, c.queue, c.routingKey)
	if err != nil {
		return nil, err
	}

	deliveries, err := ch.Consume(c.queue, "", false, true, false, false, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "can't start listening messages from amqp channel")
	}

	return deliveries, nil
}
