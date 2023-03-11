package config

import "os"

var MainConfig AppConfig

type AppConfig struct {
	ImagePath  string
	AMQPConfig AMQPConfig
}

type AMQPConfig struct {
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
		},
	}
}

func GetWithDefault(key, defaultVal string) string {
	if str := os.Getenv(key); str != "" {
		return str
	}
	return defaultVal
}
