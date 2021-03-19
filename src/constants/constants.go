package constants

import (
	"log"
	"os"
)

func init() {
	if isErr {
		log.Fatal("error constnats")
	}
}

const prefix = "SCSB" // SlackCtfScoreBot

var (
	isErr         = false
	DbUser        = NewEnv("DB_USER", true)
	DbHost        = NewEnv("DB_HOST", true)
	DbPort        = NewEnv("DB_PORT", true)
	DbPass        = NewEnv("DB_PASS", true)
	DbName        = NewEnv("DB_NAME", true)
	SlackAppToken = NewEnv("SLACK_APP_TOKEN", true)
	SlackBotToken = NewEnv("SLACK_BOT_TOKEN", true)
)

func NewEnv(key string, required bool) string {
	envkey := prefix + "_" + key
	value := os.Getenv(envkey)
	if value == "" {
		if required {
			log.Printf("%s is required\n", envkey)
			isErr = true
		} else {
			log.Printf("%s is not set\n", envkey)
		}
	}
	return value
}
