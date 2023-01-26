package grpc_client

import (
	"fmt"

	"gitlab.com/medium-project/medium_post_service/config"
	pbu "gitlab.com/medium-project/medium_post_service/genproto/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientI interface {
	UserService() pbu.UserServiceClient
}

type GrpcClient struct {
	cfg        config.Config
	connection map[string]interface{}
}

func New(cfg config.Config) (GrpcClientI, error) {
	conUserService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.UserServiceHost, cfg.UserServiceGrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("user service host: %v port: %v", cfg.UserServiceHost, cfg.UserServiceGrpcPort)
	}

	return &GrpcClient{
		cfg: cfg,
		connection: map[string]interface{}{
			"user_service": pbu.NewUserServiceClient(conUserService),
		},
	}, nil
}

func (g *GrpcClient) UserService() pbu.UserServiceClient {
	return g.connection["user_service"].(pbu.UserServiceClient)
}
