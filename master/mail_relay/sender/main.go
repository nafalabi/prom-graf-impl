package sender

import (
	"bytes"
	"fmt"
	templt "mailgun-relay/template"
	"mailgun-relay/utils"
)

var logger = utils.NewLogger(utils.INFO)

type Config map[string]EngineConfig

func SendMail(config Config, subject string, payload utils.DataPayload) error {
	for name, engineConfig := range config {
		var engine Engine
		switch engineConfig.Engine {
		case "smtp":
			engine = &SmtpEngine{}
		case "mailgun":
			engine = &MailgunEngine{}
		default:
			logger.Info(fmt.Sprintf("Can't decide sender engine for '%v'", name))
			continue
		}

		err := engine.Configure(engineConfig)
		if err != nil {
			logger.Error(fmt.Sprintf("Can't configure email sender for '%v': %v", name, err))
		}

    payload.CustomerName = engineConfig.Name

		var buff bytes.Buffer
		err = templt.GetTemplate().Execute(&buff, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("Can't render email template for '%v': %v", name, err))
		}
		mailBody := buff.String()

		err = engine.Sendmail(subject, mailBody)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	return nil
}
