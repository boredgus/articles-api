package main

import (
	"a-article/config"
	"a-article/internal/mailing"
	broker "a-article/pkg/msgbroker"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func init() {
	config.InitConfig()
}

func main() {
	var consumer broker.Consumer = broker.NewRabbitMQ()
	var mailingService = mailing.NewEmailService()
	consumer.Consume(broker.QueueParams{
		Exchange: mailing.MailingExchange,
		Queue:    "",
		Keys:     []string{mailing.MailingExchange},
		AutoAck:  true,
	}, func(d broker.Delivery) {
		var emailData mailing.Email
		if err := json.Unmarshal(d.Body(), &emailData); err != nil {
			logrus.Errorf("failed to unmarshal message: %v", err)
			return
		}
		mailingService.SendEmail(emailData)
	})
}
