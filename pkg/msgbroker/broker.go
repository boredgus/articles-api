package msgbroker

import "github.com/sirupsen/logrus"

type Delivery interface {
	ContentType() string
	Body() []byte
	Type() string
	CorrelationId() string
	ReplyTo() string
	Acknowledge(multiple bool) error
}

type Publishing struct {
	ContentType   string
	Body          []byte
	ReplyTo       string
	CorrelationId string
}

type Broker interface {
	Publisher
	Consumer
}
type QueueParams struct {
	Exchange, Queue string
	Keys            []string
	AutoAck         bool
}
type Publisher interface {
	Publish(exchange, routingKey string, msg Publishing)
}
type Consumer interface {
	Consume(queue QueueParams, processor func(d Delivery))
}

func logError(err error, msg string) {
	if err != nil {
		logrus.Panicf("%s: %s", msg, err)
	}
}
