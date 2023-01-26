package v1

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/medium-project/medium_api_gateway/api/models"
	"gitlab.com/medium-project/medium_api_gateway/config"
	grpcPkg "gitlab.com/medium-project/medium_api_gateway/pkg/grpc_client"
)

var (
	ErrWrongEmailOrPassword = errors.New("wrong email or password")
	ErrUserNotVerifid       = errors.New("user not verified")
	ErrEmailExists          = errors.New("email is already exists")
	ErrIncorrectCode        = errors.New("incorrect verification code")
	ErrCodeExpired          = errors.New("verification is expired")
	ErrForbidden            = errors.New("forbidden")
	ErrNotAllowed           = errors.New("method not allowed")
	ErrWeakPassword         = errors.New("password must contain at least one small letter, one capital letter, one number, one symbol")
)

const (
	RegisterCodeKey   = "register_code_"
	ForgotPasswordKey = "forgot_password_key_"
)

type handlerV1 struct {
	cfg        *config.Config
	grpcClient grpcPkg.GrpcClientI
	logger     *logrus.Logger
}

type HandlerV1Options struct {
	Cfg        *config.Config
	GrpcClient grpcPkg.GrpcClientI
	Logger     *logrus.Logger
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		cfg:        options.Cfg,
		grpcClient: options.GrpcClient,
		logger:     options.Logger,
	}
}

func errResponse(err error) *models.ResponseError {
	return &models.ResponseError{
		Error: err.Error(),
	}
}

func validateGetAllParams(ctx *gin.Context) (*models.GetAllParams, error) {
	var (
		limit int64 = 10
		page  int64 = 1
		err   error
	)
	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &models.GetAllParams{
		Limit:  limit,
		Page:   page,
		Search: ctx.Query("search"),
	}, nil
}

func validateGetAllPostsParams(ctx *gin.Context) (*models.GetAllPostsParams, error) {
	var (
		limit              int64 = 10
		page               int64 = 1
		err                error
		userId, categoryId int64
		sortByDate         string
	)
	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("user_id") != "" {
		userId, err = strconv.ParseInt(ctx.Query("user_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("category_id") != "" {
		categoryId, err = strconv.ParseInt(ctx.Query("category_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if ctx.Query("sort") != "" &&
		(ctx.Query("sort") == "desc" || ctx.Query("sort") == "asc") {
		sortByDate = ctx.Query("sort")
	}

	return &models.GetAllPostsParams{
		Limit:      limit,
		Page:       page,
		Search:     ctx.Query("search"),
		UserID:     userId,
		CategoryID: categoryId,
		SortByDate: sortByDate,
	}, nil
}

func validateGetAllCommentsParams(ctx *gin.Context) (*models.GetAllCommentsParams, error) {
	var (
		limit  int64 = 10
		page   int64 = 1
		err    error
		postId int64
	)

	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("post_id") != "" {
		postId, err = strconv.ParseInt(ctx.Query("post_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &models.GetAllCommentsParams{
		Limit:  limit,
		Page:   page,
		SortBy: ctx.Query("sort_by"),
		PostID: postId,
	}, nil
}
