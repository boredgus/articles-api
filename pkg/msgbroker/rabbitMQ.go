package msgbroker

import (
	"a-article/config"
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMQDelivery struct {
	amqp.Delivery
}

func (d RabbitMQDelivery) Body() []byte {
	return d.Delivery.Body
}
func (d RabbitMQDelivery) Type() string {
	return d.Delivery.Type
}
func (d RabbitMQDelivery) CorrelationId() string {
	return d.Delivery.CorrelationId
}
func (d RabbitMQDelivery) Acknowledge(multiple bool) error {
	return d.Delivery.Ack(multiple)
}
func (d RabbitMQDelivery) ReplyTo() string {
	return d.Delivery.ReplyTo
}
func (d RabbitMQDelivery) ContentType() string {
	return d.Delivery.ContentType
}

type RabbitMQ struct {
	dsn string
}

func NewRabbitMQ() Broker {
	env := config.GetConfig()
	return &RabbitMQ{
		dsn: fmt.Sprintf("amqp://%s:%s@%s:%s/",
			env.RabbitMQUser, env.RabbitMQPass, env.RabbitMQContainer, env.RabbitMQPort),
	}
}
func (b *RabbitMQ) connect(processor func(ch *amqp.Channel)) {
	conn, err := amqp.Dial(b.dsn)
	logError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	logError(err, "failed to open channel")
	defer ch.Close()
	processor(ch)
}

func (b *RabbitMQ) Publish(exchange, routingKey string, msg Publishing) {
	b.connect(func(ch *amqp.Channel) {
		if exchange != "" {
			if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
				logrus.Fatalf("failed to create exchange %s: %v", exchange, err)
			}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := ch.PublishWithContext(ctx,
			exchange,   // exchange
			routingKey, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType:   msg.ContentType,
				Body:          msg.Body,
				ReplyTo:       msg.ReplyTo,
				CorrelationId: msg.CorrelationId,
			}); err != nil {
			logrus.Errorf("failed to publish a message: %v", err)
		}
	})
}

func (b *RabbitMQ) Consume(params QueueParams, processor func(d Delivery)) {
	b.connect(func(ch *amqp.Channel) {
		if params.Exchange != "" {
			if err := ch.ExchangeDeclare(params.Exchange, "topic", true, false, false, false, nil); err != nil {
				logrus.Fatalf("failed to create exchange %s: %v", params.Exchange, err)
			}
		}
		q, err := ch.QueueDeclare(params.Queue, true, false, false, false, nil)
		if err != nil {
			logrus.Fatalf("failed to create queue: %v", err)
		}
		for _, key := range params.Keys {
			if err := ch.QueueBind(q.Name, key, params.Exchange, false, nil); err != nil {
				logrus.Fatalf("failed to bind queue %s (routing key: %s, exchange: %s): %v", q.Name, key, params.Exchange, err)
			}
		}

		msgs, err := ch.Consume(
			q.Name,         // queue
			"",             // consumer
			params.AutoAck, // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		logError(err, "failed to register a consumer")

		var forever chan struct{}
		go func() {
			for d := range msgs {
				processor(RabbitMQDelivery{d})
			}
		}()
		<-forever
	})
}
