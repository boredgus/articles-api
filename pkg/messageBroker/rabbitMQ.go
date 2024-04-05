package messagbroker

import (
	"a-article/config"
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQDelivery amqp.Delivery

func (d RabbitMQDelivery) Data() []byte {
	return d.Body
}
func (d RabbitMQDelivery) MessageType() string {
	return d.Type
}

type RabbitMQ struct {
	conn *amqp.Connection
}

func NewRabbitMQ() Broker {
	config := config.GetConfig()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.RabbitMQUser, config.RabbitMQPass, config.RabbitMQContainer, config.RabbitMQPort))
	logError(err, "Failed to connect to RabbitMQ")
	return &RabbitMQ{conn: conn}
}

func (b *RabbitMQ) declareQueue(name string, ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (b *RabbitMQ) Publish(message []byte, queue string) {
	ch, err := b.conn.Channel()
	logError(err, "failed to open channel on "+b.conn.Config.Vhost)
	defer ch.Close()

	q, err := b.declareQueue(queue, ch)
	logError(err, "failed to declare queue '"+queue+"'")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "json/application",
			Body:        message,
		})
	logError(err, "failed to publish a message")
}

func (b *RabbitMQ) Consume(queue string, processor func(d Delivery)) {
	ch, err := b.conn.Channel()
	logError(err, "failed to open channel on "+b.conn.Config.Vhost)
	defer ch.Close()

	q, err := b.declareQueue(queue, ch)
	logError(err, "failed to declare queue '"+queue+"'")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	logError(err, "failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			processor(RabbitMQDelivery(d))
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
