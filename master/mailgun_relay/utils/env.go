package utils

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvVar struct {
	Domain string
	ApiKey string
	Sender string
}

func (env *EnvVar) loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	env.Domain = os.Getenv("MAILGUN_DOMAIN")
	env.ApiKey = os.Getenv("MAILGUN_API_KEY")
	env.Sender = os.Getenv("SENDER_MAIL_ADDRESS")
}

func NewEnv() EnvVar {
	env := EnvVar{}
	env.loadEnv()
	return env
}
