package sender

import (
	"fmt"
	"mailgun-relay/utils"
)

var logger = utils.NewLogger(utils.INFO)

type Config map[string]EngineConfig

func SendMail(config Config, subject string, body string) error {
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

    err = engine.Sendmail(subject, body)
		if err != nil {
      logger.Error(err.Error())
		}
	}
  return nil
}
