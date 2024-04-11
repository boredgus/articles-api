package mailing

import (
	broker "a-article/pkg/msgbroker"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Mailman interface {
	ConfirmSignupEmail(email, passcode string)
	WelcomeEmail(email string)
}

const MailingExchange = "mailing"

type mailing struct {
	broker broker.Publisher
}

func NewMailman(broker broker.Publisher) Mailman {
	return &mailing{broker: broker}
}

func (m *mailing) publishMessage(data []byte) {
	m.broker.Publish(MailingExchange, MailingExchange, broker.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
}

func (m *mailing) ConfirmSignupEmail(email, passcode string) {
	msg, err := json.Marshal(Email{
		Receiver: email,
		Template: string(ConfirmSignupEmail),
		Params:   map[string]any{"passcode": passcode},
	})
	if err != nil {
		logrus.Errorf("failed to marshal email message data: %v", err)
		return
	}
	m.publishMessage(msg)
}

func (m *mailing) WelcomeEmail(email string) {
	msg, err := json.Marshal(Email{
		Receiver: email,
		Template: string(WelcomeEmail),
	})
	if err != nil {
		logrus.Errorf("failed to marshal email message data: %v", err)
		return
	}
	m.publishMessage(msg)
}
