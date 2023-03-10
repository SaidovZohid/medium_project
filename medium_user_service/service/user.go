package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	pb "gitlab.com/medium-project/medium_user_service/genproto/user_service"
	"gitlab.com/medium-project/medium_user_service/pkg/utils"
	"gitlab.com/medium-project/medium_user_service/storage"
	"gitlab.com/medium-project/medium_user_service/storage/repo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
	logger   *logrus.Logger
}

func NewUserService(strg storage.StorageI, inMemory storage.InMemoryStorageI, log *logrus.Logger) *UserService {
	return &UserService{
		storage:  strg,
		inMemory: inMemory,
		logger:   log,
	}
}

func (s *UserService) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.WithError(err).Error("failed to hash password")
		return nil, status.Errorf(codes.Internal, "internatl server error: %v", err)
	}
	user, err := s.storage.User().Create(&repo.User{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Gender:          req.Gender,
		Password:        hashedPassword,
		Username:        req.Username,
		ProfileImageUrl: req.ProfileImageUrl,
		Type:            req.Type,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create user")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.User{
		Id:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Gender:          user.Gender,
		Password:        user.Password,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) Get(ctx context.Context, req *pb.IdRequest) (*pb.User, error) {
	user, err := s.storage.User().Get(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user in get func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.User{
		Id:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Gender:          user.Gender,
		Password:        user.Password,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) GetByEmail(ctx context.Context, req *pb.GetByEmailRequest) (*pb.User, error) {
	user, err := s.storage.User().GetByEmail(req.Email)
	if err != nil {
		s.logger.WithError(err).Error("failed to getbyemail user in get func")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.User{
		Id:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Gender:          user.Gender,
		Password:        user.Password,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) Update(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().Update(&repo.User{
		ID:              req.Id,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Username:        req.Username,
		Gender:          req.Gender,
		Email:           req.Email,
		PhoneNumber:     req.PhoneNumber,
		ProfileImageUrl: req.ProfileImageUrl,
		Type:            req.Type,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to update user in update func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.User{
		Id:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Gender:          user.Gender,
		Password:        user.Password,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) Delete(ctx context.Context, req *pb.IdRequest) (*emptypb.Empty, error) {
	err := s.storage.User().Delete(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete user in delete func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) GetAll(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.storage.User().GetAll(&repo.GetAllUserParams{
		Limit:  req.Limit,
		Page:   req.Page,
		Search: req.Search,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to get all user in getall func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return parseUserResponse(users)
}

func parseUserResponse(users *repo.GetAllUsersResult) (*pb.GetAllUsersResponse, error) {
	var res pb.GetAllUsersResponse
	res.Count = users.Count
	for _, user := range users.Users {
		u := pb.User{
			Id:              user.ID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Gender:          user.Gender,
			Email:           user.Email,
			PhoneNumber:     user.PhoneNumber,
			Username:        user.Username,
			ProfileImageUrl: user.ProfileImageUrl,
			Type:            user.Type,
			CreatedAt:       user.CreatedAt.Format(time.RFC3339),
		}
		res.Users = append(res.Users, &u)
	}

	return &res, nil
}
