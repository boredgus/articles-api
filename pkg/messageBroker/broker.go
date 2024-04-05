package messagbroker

import "github.com/sirupsen/logrus"

type Delivery interface {
	Data() []byte
	MessageType() string
}

type Broker interface {
	Publisher
	Consumer
}
type Publisher interface {
	Publish(message []byte, queueName string)
}
type Consumer interface {
	Consume(queueName string, processor func(d Delivery))
}

func logError(err error, msg string) {
	if err != nil {
		logrus.Panicf("%s: %s", msg, err)
	}
}
