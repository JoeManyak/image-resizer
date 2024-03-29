package config

import (
	"os"
	"strconv"
)

var MainConfig AppConfig

type AppConfig struct {
	ImagePath  string
	AMQPConfig AMQPConfig
}

type AMQPConfig struct {
	Retries   int
	Timeout   int
	QueueName string
	User      string
	Pass      string
	URL       string
	Port      string
}

func Setup() {
	MainConfig = AppConfig{
		ImagePath: GetWithDefault(ImagePath, "./img"),
		AMQPConfig: AMQPConfig{
			QueueName: GetWithDefault(AMQPQueueName, "main"),
			User:      GetWithDefault(AMQPUser, "guest"),
			Pass:      GetWithDefault(AMQPPass, "guest"),
			URL:       GetWithDefault(AMQPUrl, "localhost"),
			Port:      GetWithDefault(AMQPPort, "5672"),
			Retries:   GetWithDefaultNumber(AMQPRetries, 5),
			Timeout:   GetWithDefaultNumber(AMQPTimeout, 3),
		},
	}
}

func GetWithDefault(key, defaultVal string) string {
	if str := os.Getenv(key); str != "" {
		return str
	}
	return defaultVal
}

func GetWithDefaultNumber(key string, defaultVal int) int {
	if str := os.Getenv(key); str != "" {
		val, err := strconv.Atoi(str)
		if err != nil {
			return defaultVal
		}
		return val
	}
	return defaultVal
}
