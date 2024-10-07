package sender

type SmtpConfig struct {
	Host     string `json:"HOST"`
	Port     string `json:"PORT"`
	Username string `json:"USERNAME"`
	Password string `json:"PASSWORD"`
}

type MailgunConfig struct {
	MailgunDomain string `json:"MAILGUN_DOMAIN"`
	MailgunApiKey string `json:"MAILGUN_APIKEY"`
}

type EngineConfig struct {
	Name          string         `json:"name"`
	Engine        string         `json:"engine"` // smtp or mailgun
	Sender        string         `json:"sender"`
	Receivers     []string       `json:"receivers"`
	SmtpConfig    *SmtpConfig    `json:"smtp_config"`
	MailgunConfig *MailgunConfig `json:"mailgun_config"`
}

type Engine interface {
	Configure(config EngineConfig) error
	Sendmail(subject string, body string) error
}
