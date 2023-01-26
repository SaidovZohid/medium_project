package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "gitlab.com/medium-project/medium_post_service/genproto/post_service"
	"gitlab.com/medium-project/medium_post_service/genproto/user_service"
	grpcPkg "gitlab.com/medium-project/medium_post_service/pkg/grpc_client"
	"gitlab.com/medium-project/medium_post_service/storage"
	"gitlab.com/medium-project/medium_post_service/storage/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
	storage    storage.StorageI
	grpcClient grpcPkg.GrpcClientI
	logger     *logrus.Logger
}

func NewCommentService(strg *storage.StorageI, grpc grpcPkg.GrpcClientI, log *logrus.Logger) *CommentService {
	return &CommentService{
		storage:    *strg,
		grpcClient: grpc,
		logger:     log,
	}
}

func (s *CommentService) Create(ctx context.Context, req *pb.Comment) (*pb.Comment, error) {
	comment, err := s.storage.Comment().Create(&repo.Comment{
		PostID:      req.PostId,
		UserID:      req.UserId,
		Description: req.Description,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create comment")
		return nil, status.Errorf(codes.Internal, "internal server error in create comment service: %v", err)
	}

	user, err := s.grpcClient.UserService().Get(context.Background(), &user_service.IdRequest{
		Id: req.UserId,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to get user info in create comment func")
		return nil, status.Errorf(codes.Internal, "internal server error in while getting user info comment service: %v", err)
	}

	return &pb.Comment{
		Id:          comment.ID,
		PostId:      comment.PostID,
		UserId:      comment.UserID,
		Description: comment.Description,
		CreatedAt:   comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   comment.UpdatedAt.Format(time.RFC3339),
		User: &pb.CommentUser{
			Id:              user.Id,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			ProfileImageUrl: user.ProfileImageUrl,
		},
	}, nil
}
func (s *CommentService) Delete(ctx context.Context, req *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
	err := s.storage.Comment().Delete(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete comment")
		return nil, status.Errorf(codes.Internal, "internale server while deleting comment: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *CommentService) GetAll(ctx context.Context, req *pb.GetAllCommentsParamsReq) (*pb.GetAllCommentsResponse, error) {
	comments, err := s.storage.Comment().GetAll(&repo.GetCommentsParams{
		Limit:  req.Limit,
		Page:   req.Page,
		PostID: req.PostId,
		SortBy: req.SortBy,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to get all comments")
		return nil, status.Errorf(codes.Internal, "internal server error while updating comment: %v", err)
	}

	res := pb.GetAllCommentsResponse{
		Comments: make([]*pb.Comment, 0),
		Count:    comments.Count,
	}

	for _, comment := range comments.Comments {
		user, err := s.grpcClient.UserService().Get(context.Background(), &user_service.IdRequest{Id: comment.UserID})
		if err != nil {
			s.logger.WithError(err).Error("failed to get user info from user service")
			return nil, status.Errorf(codes.Internal, "internal server error while updating comment: %v", err)
		}
		res.Comments = append(res.Comments, parseComment(user, comment))
	}

	return &res, nil
}

func parseComment(user *user_service.User, comment *repo.Comment) *pb.Comment {
	return &pb.Comment{
		Id:          comment.ID,
		PostId:      comment.PostID,
		UserId:      comment.UserID,
		Description: comment.Description,
		CreatedAt:   comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   comment.UpdatedAt.Format(time.RFC3339),
		User: &pb.CommentUser{
			Id:              user.Id,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			ProfileImageUrl: user.ProfileImageUrl,
		},
	}
}

func (s *CommentService) Update(ctx context.Context, req *pb.Comment) (*pb.Comment, error) {
	comment, err := s.storage.Comment().Update(&repo.Comment{
		Description: req.Description,
		ID:          req.Id,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to update comment")
		return nil, status.Errorf(codes.Internal, "internal server error while updating comment: %v", err)
	}
	user, err := s.grpcClient.UserService().Get(context.Background(), &user_service.IdRequest{Id: comment.UserID})
	if err != nil {
		s.logger.WithError(err).Error("failed to get user info from user service in update func")
		return nil, status.Errorf(codes.Internal, "internal server error while updating comment: %v", err)
	}
	return &pb.Comment{
		Id:          comment.ID,
		PostId:      comment.PostID,
		UserId:      comment.UserID,
		Description: comment.Description,
		CreatedAt:   comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   comment.UpdatedAt.Format(time.RFC3339),
		User: &pb.CommentUser{
			Id:              user.Id,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			ProfileImageUrl: user.ProfileImageUrl,
		},
	}, nil
}
