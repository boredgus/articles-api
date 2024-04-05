package mailing

import (
	broker "a-article/pkg/messageBroker"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Mailman interface {
	ConfirmSignupEmail(email, passcode string)
	WelcomeEmail(email string)
}

const MailingQueue = "mailing"

type mailing struct {
	broker broker.Publisher
}

func NewMailman() Mailman {
	return &mailing{broker: broker.NewRabbitMQ()}
}

func (m *mailing) ConfirmSignupEmail(email, passcode string) {
	msg, err := json.Marshal(Email{
		Receiver: email,
		Template: string(ConfirmSignupEmail),
		Params:   map[string]any{"passcode": passcode},
	})
	if err != nil {
		logrus.Errorf("failed to marshal email message data: %v", err)
	}
	m.broker.Publish(msg, MailingQueue)
}

func (m *mailing) WelcomeEmail(email string) {
	msg, err := json.Marshal(Email{
		Receiver: email,
		Template: string(WelcomeEmail),
	})
	if err != nil {
		logrus.Errorf("failed to marshal email message data: %v", err)
	}
	m.broker.Publish(msg, MailingQueue)
}
