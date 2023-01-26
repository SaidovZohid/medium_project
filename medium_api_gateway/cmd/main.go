package main

import (
	"log"

	"gitlab.com/medium-project/medium_api_gateway/api"
	_ "gitlab.com/medium-project/medium_api_gateway/api/docs"
	"gitlab.com/medium-project/medium_api_gateway/config"
	grpcPkg "gitlab.com/medium-project/medium_api_gateway/pkg/grpc_client"
	"gitlab.com/medium-project/medium_api_gateway/pkg/logger"
)

func main() {
	cfg := config.Load(".")

	grpcConn, err := grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to get grpc connettion: %v", err)
	}

	logger := logger.New()

	apiServer := api.New(&api.RoutetOptions{
		Cfg:        &cfg,
		GrpcClient: grpcConn,
		Logger:     logger,
	})

	err = apiServer.Run(cfg.HttpPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
