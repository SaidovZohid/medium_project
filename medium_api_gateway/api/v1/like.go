package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/medium-project/medium_api_gateway/api/models"
	pb "gitlab.com/medium-project/medium_api_gateway/genproto/post_service"
)

// @Security ApiKeyAuth
// @Router /likes [post]
// @Summary Create Or Update like
// @Description Create Or Update like
// @Tags like
// @Accept json
// @Produce json
// @Param like body models.CreateOrUpdateLikeRequest true "like"
// @Success 201 {object} models.Like
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) CreateOrUpdateLike(ctx *gin.Context) {
	var (
		req models.CreateOrUpdateLikeRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		h.logger.WithError(err).Error("failed to get auth payload")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	resp, err := h.grpcClient.LikeService().CreateOrUpdate(context.Background(), &pb.Like{
		Status: req.Status,
		PostId: req.PostID,
		UserId: payload.UserID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.Like{
		ID:     resp.Id,
		PostID: resp.PostId,
		UserID: resp.UserId,
		Status: resp.Status,
	})
}

// @Security ApiKeyAuth
// @Router /likes/user-post/{id} [get]
// @Summary Get like by giving to query post_id
// @Description Get like by giving to query post_id
// @Tags like
// @Accept json
// @Produce json
// @Param id path int true "post_id"
// @Success 201 {object} models.Like
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetLike(ctx *gin.Context) {
	postID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		h.logger.WithError(err).Error("failed to get auth payload")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	resp, err := h.grpcClient.LikeService().Get(context.Background(), &pb.GetLikeRequest{
		UserId: payload.UserID,
		PostId: int64(postID),
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to get like")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.Like{
		ID:     resp.Id,
		PostID: resp.PostId,
		UserID: resp.UserId,
		Status: resp.Status,
	})
}

// @Router /likes/user-post-likes/{id} [get]
// @Summary Get likes and dislike count by giving to query post_id
// @Description Get likes and dislikes count by giving to query post_id
// @Tags like
// @Accept json
// @Produce json
// @Param id path int true "post_id"
// @Success 201 {object} models.LikesAndDislikesCount
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetLikesAndDislikesCount(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		h.logger.WithError(err).Error("failed to parse id or bad request")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	count, err := h.grpcClient.LikeService().GetLikesDislikesCount(context.Background(), &pb.GetLikesRequest{
		PostId: id,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.LikesAndDislikesCount{
		Likes:    count.Likes,
		Dislikes: count.Dislikes,
	})
}
