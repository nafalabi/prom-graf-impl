package sender

import (
	"bytes"
	"fmt"
	templt "mailgun-relay/template"
	"mailgun-relay/utils"
)

var logger = utils.NewLogger(utils.INFO)

func SendMail(senderConfig SenderConfig, subject string, payload utils.DataPayload) error {
	var engine Engine

	switch senderConfig.Engine {
	case "smtp":
		engine = &SmtpEngine{}
	case "mailgun":
		engine = &MailgunEngine{}
	default:
		return fmt.Errorf("Can't decide sender engine for '%v'", senderConfig.Name)
	}

	err := engine.Configure(senderConfig)
	if err != nil {
		return fmt.Errorf("Can't configure email sender for '%v': %v", senderConfig.Name, err)
	}

	payload.CustomerName = senderConfig.Name

	var buff bytes.Buffer
	err = templt.GetTemplate().Execute(&buff, payload)
	if err != nil {
		return fmt.Errorf("Can't render email template for '%v': %v", senderConfig.Name, err)
	}
	mailBody := buff.String()

	err = engine.Sendmail(subject, mailBody)
	if err != nil {
		return err
	}

	return nil
}
