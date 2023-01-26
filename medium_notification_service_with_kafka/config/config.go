package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort string
	KafkaUrl string
	Smtp     Smtp
}

type Smtp struct {
	Sender   string
	Password string
}

func Load(path string) Config {
	godotenv.Load(path + "/.env")

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		GrpcPort: conf.GetString("GRPC_PORT"),
		KafkaUrl: conf.GetString("KAFKA_URL"),
		Smtp: Smtp{
			Sender:   conf.GetString("SMTP_SENDER"),
			Password: conf.GetString("SMTP_PASSWORD"),
		},
	}
	return cfg
}
