package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "gitlab.com/medium-project/medium_post_service/genproto/post_service"
	"gitlab.com/medium-project/medium_post_service/storage"
	"gitlab.com/medium-project/medium_post_service/storage/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PostService struct {
	pb.UnimplementedPostServiceServer
	storage storage.StorageI
	logger  *logrus.Logger
}

func NewPostService(strg *storage.StorageI, log *logrus.Logger) *PostService {
	return &PostService{
		storage: *strg,
		logger:  log,
	}
}

func (s *PostService) Create(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	post, err := s.storage.Post().Create(&repo.Post{
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		UserID:      req.UserId,
		CategoryID:  req.CategoryId,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create post in create func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return parsePost(post), nil
}

func (s *PostService) Get(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	post, err := s.storage.Post().Get(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to get post in get func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return parsePost(post), nil
}

func (s *PostService) Update(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	post, err := s.storage.Post().Update(&repo.Post{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		UserID:      req.UserId,
		CategoryID:  req.CategoryId,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to update post in update func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return parsePost(post), nil
}

func (s *PostService) Delete(ctx context.Context, req *pb.GetPostRequest) (*emptypb.Empty, error) {
	err := s.storage.Post().Delete(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete post in delete func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *PostService) GetAll(ctx context.Context, req *pb.GetPostsParamsReq) (*pb.GetAllPostResponse, error) {
	posts, err := s.storage.Post().GetAll(&repo.GetPostsParams{
		Limit:      req.Limit,
		Page:       req.Page,
		Search:     req.Search,
		UserID:     req.UserId,
		CategoryID: req.CategoryId,
		SortByDate: req.SortByDate,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to get all post in getall func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	res := pb.GetAllPostResponse{
		Posts: make([]*pb.Post, 0),
		Count: posts.Count,
	}
	for _, post := range posts.Posts {
		p := parsePost(post)
		res.Posts = append(res.Posts, p)
	}
	return &res, nil
}

func parsePost(req *repo.Post) *pb.Post {
	return &pb.Post{
		Id:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		UserId:      req.UserID,
		CategoryId:  req.CategoryID,
		CreatedAt:   req.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   req.UpdatedAt.Format(time.RFC3339),
		ViewsCount:  int64(req.ViewsCount),
	}
}
