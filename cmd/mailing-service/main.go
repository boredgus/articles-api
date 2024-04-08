package main

import (
	"a-article/config"
	"a-article/internal/mailing"
	broker "a-article/pkg/messageBroker"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func init() {
	config.InitConfig()
}

func main() {
	var consumer broker.Consumer = broker.NewRabbitMQ()
	var mailingService = mailing.NewEmailService()
	consumer.Consume(mailing.MailingQueue, func(d broker.Delivery) {
		logrus.Infof("%+v", string(d.Data()))
		var emailData mailing.Email
		if err := json.Unmarshal(d.Data(), &emailData); err != nil {
			logrus.Errorf("failed to unmarshal message: %v", err)
			return
		}
		mailingService.SendEmail(emailData)
	})
}
