package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HttpPort            string
	UserServiceGrpcPort string
	UserServiceHost     string
	PostServiceGrpcPort string
	PostServiceHost     string
}

func Load(path string) Config {
	godotenv.Load(path + "/.env")

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		HttpPort:            conf.GetString("HTTP_PORT"),
		UserServiceGrpcPort: conf.GetString("USER_SERVICE_GRPC_PORT"),
		UserServiceHost:     conf.GetString("USER_SERVICE_HOST"),
		PostServiceGrpcPort: conf.GetString("POST_SERVICE_GRPC_PORT"),
		PostServiceHost:     conf.GetString("POST_SERVICE_HOST"),
	}
	return cfg
}
