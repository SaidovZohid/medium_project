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

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	storage storage.StorageI
	logger  *logrus.Logger
}

func NewCategoryService(strg *storage.StorageI, log *logrus.Logger) *CategoryService {
	return &CategoryService{
		storage: *strg,
		logger:  log,
	}
}

func (s *CategoryService) Create(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	category, err := s.storage.Category().Create(&repo.Category{
		Title: req.Title,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create new category in create func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return parseCategory(category), nil
}

func (s *CategoryService) Get(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	category, err := s.storage.Category().Get(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to get category in get func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return parseCategory(category), nil
}

func (s *CategoryService) Update(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	category, err := s.storage.Category().Update(&repo.Category{
		ID:    req.Id,
		Title: req.Title,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to update category in update func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return parseCategory(category), nil
}

func (s *CategoryService) Delete(ctx context.Context, req *pb.GetCategoryRequest) (*emptypb.Empty, error) {
	err := s.storage.Category().Delete(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete category in delete func ")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *CategoryService) GetAll(ctx context.Context, req *pb.GetAllCategoryParamsReq) (*pb.GetAllCategoryResponse, error) {
	categories, err := s.storage.Category().GetAll(&repo.GetAllCategoryParams{
		Limit:  req.Limit,
		Page:   req.Page,
		Search: req.Search,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to get all categories in getall func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	res := pb.GetAllCategoryResponse{
		Categories: make([]*pb.Category, 0),
		Count:      int64(categories.Count),
	}
	for _, category := range categories.Categories {
		p := parseCategory(category)
		res.Categories = append(res.Categories, p)
	}
	return &res, nil
}

func parseCategory(req *repo.Category) *pb.Category {
	return &pb.Category{
		Id:        req.ID,
		Title:     req.Title,
		CreatedAt: req.CreatedAt.Format(time.RFC3339),
	}
}
