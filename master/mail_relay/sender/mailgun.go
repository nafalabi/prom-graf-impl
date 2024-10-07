package sender

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunEngine struct {
	Domain    string
	ApiKey    string
	Sender    string
	Receivers []string
}

func (m *MailgunEngine) Configure(config EngineConfig) error {
  if (config.MailgunConfig == nil) {
    return errors.New("mailgun config was not provided")
  }

  mailgunConfig := config.MailgunConfig

  m.Domain = mailgunConfig.MailgunDomain
  m.ApiKey = mailgunConfig.MailgunApiKey
  m.Sender = config.Sender
  m.Receivers = config.Receivers
  
	return nil
}

func (m *MailgunEngine) Sendmail(subject string, body string) error {
	mg := mailgun.NewMailgun(m.Domain, m.ApiKey)

	message := mg.NewMessage(m.Sender, subject, "")
	for _, recipient := range m.Receivers {
		err := message.AddRecipient(recipient)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}

	message.SetHtml(body)
	message.AddInline("./bawana.png")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}

	logger.Info(fmt.Sprintf("Mail sent successfully to %v", m.Receivers))
	return nil
}
