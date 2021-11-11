package rmq

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type publisher struct {
	log              *zap.Logger
	connectionString string
	channel          *amqp.Channel
	closeCh          chan *amqp.Error

	queue      string
	routingKey string
	exchange   string
}

// NewPublisher create new publisher
func (c *rmq) NewPublisher(exchange, queue, routingKey string) (*publisher, error) {
	publisher := &publisher{
		log:              c.log.Named(fmt.Sprintf("publisher(%s::%s)", queue, routingKey)),
		connectionString: c.connectionString,
		exchange:         exchange,
		routingKey:       routingKey,
		queue:            queue,
	}

	if err := publisher.connect(); err != nil {
		c.log.Error("publisher connection",
			zap.String("connection_string", c.connectionString),
			zap.String("exchange", exchange),
			zap.String("queue", queue),
			zap.String("routing_key", routingKey),
		)
		return nil, errors.Wrap(err, "publisher connect")
	}

	return publisher, nil
}

func (p *publisher) connect() error {
	if err := p.tryConnect(); err != nil {
		return errors.Wrapf(err, "connect to amqp failed")
	}

	go p.keepConnection()

	return nil
}

func (p *publisher) tryConnect() error {
	ch, err := connect(p.connectionString, p.exchange, p.queue, p.routingKey)
	if err != nil {
		return errors.Wrapf(err, "can't open connection")
	}

	p.closeCh = make(chan *amqp.Error)
	ch.NotifyClose(p.closeCh)
	p.channel = ch

	return nil
}

func (p *publisher) keepConnection() {
	go func() {
		for {
			closeErr := <-p.closeCh
			p.log.Error("connection was closed", zap.Error(closeErr))

			for {
				err := p.tryConnect()
				if err != nil {
					p.log.Error("try to connect", zap.Error(err), zap.Int("delay_sec", int(reconnectDelay.Seconds())))
					time.Sleep(reconnectDelay)
					continue
				}

				break
			}

			p.log.Info("reconnected")
		}
	}()
}

func (p *publisher) Publish(message []byte) error {
	err := p.channel.Publish(
		p.exchange,
		p.routingKey,
		true,
		false,
		amqp.Publishing{
			Body:         message,
			DeliveryMode: amqp.Persistent,
			Priority:     1,
		},
	)

	if err != nil {
		return errors.Wrap(err, "publish message")
	}

	return nil
}
