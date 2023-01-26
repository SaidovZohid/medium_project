package main

import (
	"fmt"
	"log"
	"net"

	"gitlab.com/medium-project/medium_post_service/config"
	pb "gitlab.com/medium-project/medium_post_service/genproto/post_service"
	grpcPkg "gitlab.com/medium-project/medium_post_service/pkg/grpc_client"
	"gitlab.com/medium-project/medium_post_service/pkg/logger"
	"gitlab.com/medium-project/medium_post_service/service"
	"gitlab.com/medium-project/medium_post_service/storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	strg := storage.NewStoragePg(psqlConn)

	log := logger.New()

	grpcConn, err := grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	postService := service.NewPostService(&strg, log)
	categoryService := service.NewCategoryService(&strg, log)
	likeService := service.NewLikeService(&strg, log)
	commentService := service.NewCommentService(&strg, grpcConn, log)

	listen, err := net.Listen("tcp", cfg.GrpcPort)

	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, postService)
	pb.RegisterCategoryServiceServer(s, categoryService)
	pb.RegisterLikeServiceServer(s, likeService)
	pb.RegisterCommentServiceServer(s, commentService)

	reflection.Register(s)

	log.Info("gRPC server started port in: ", cfg.GrpcPort)
	if s.Serve(listen); err != nil {
		log.Fatalf("Error while listening: %v", err)
	}
}
