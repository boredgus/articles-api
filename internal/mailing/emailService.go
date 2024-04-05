package mailing

import (
	"a-article/config"
	"context"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/sirupsen/logrus"
)

type EmailTemplate = string

const (
	ConfirmSignupEmail EmailTemplate = "confirm_signup"
	WelcomeEmail       EmailTemplate = "welcome"
)

var templateNameToId = map[string]int64{
	ConfirmSignupEmail: 2,
	WelcomeEmail:       6,
}

type Email struct {
	Receiver string         `json:"email"`
	Template string         `json:"template"`
	Params   map[string]any `json:"params"`
}

type EmailService interface {
	SendEmail(email Email)
}

func NewEmailService() EmailService {
	cfg := config.GetConfig()
	apiConfig := brevo.NewConfiguration()
	apiConfig.AddDefaultHeader("api-key", cfg.BrevoAPIKey)
	apiConfig.AddDefaultHeader("From", cfg.SMTPUsername)
	return &emailService{
		apiClient: brevo.NewAPIClient(apiConfig),
	}
}

type emailService struct {
	apiClient *brevo.APIClient
}

func (s *emailService) SendEmail(email Email) {
	if _, _, err := s.apiClient.TransactionalEmailsApi.SendTransacEmail(context.Background(), brevo.SendSmtpEmail{
		To:         []brevo.SendSmtpEmailTo{{Email: email.Receiver}},
		TemplateId: templateNameToId[email.Template],
		Params:     email.Params,
	}); err != nil {
		logrus.Errorf("failed to send email: %v", err)
	}
}
