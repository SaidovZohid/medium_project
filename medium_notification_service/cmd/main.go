package main

import (
	"log"
	"net"

	_ "github.com/lib/pq"
	"gitlab.com/medium-project/medium_notification_service/config"
	pb "gitlab.com/medium-project/medium_notification_service/genproto/notification_service"
	"gitlab.com/medium-project/medium_notification_service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load(".")

	notificationService := service.NewNotificationService(&cfg)

	listen, err := net.Listen("tcp", cfg.GrpcPort)

	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, notificationService)
	reflection.Register(s)

	log.Println("gRPC server started port in: ", cfg.GrpcPort)
	if s.Serve(listen); err != nil {
		log.Fatalf("Error while listening: %v", err)
	}
}
