package v1_test

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/medium-project/medium_api_gateway/api"
	"gitlab.com/medium-project/medium_api_gateway/config"
	grpcPkg "gitlab.com/medium-project/medium_api_gateway/pkg/grpc_client"
	"gitlab.com/medium-project/medium_api_gateway/pkg/logger"
)

var (
	router *gin.Engine
	grpcConn grpcPkg.GrpcClientI
)

func TestMain(m *testing.M) {
	var err  error
	cfg := config.Load("./../..")

	logrus := logger.New()

	grpcConn, err = grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to get grpc connection %v", err)
	}

	ginEngine := api.New(&api.RoutetOptions{
		Cfg:        &cfg,
		GrpcClient: grpcConn,
		Logger:     logrus,
	})

	router = ginEngine

	os.Exit(m.Run())
}
