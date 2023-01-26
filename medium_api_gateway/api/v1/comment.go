package v1

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/medium-project/medium_api_gateway/api/models"

	pb "gitlab.com/medium-project/medium_api_gateway/genproto/post_service"
)

// @Security ApiKeyAuth
// @Router /comments [post]
// @Summary Create a comment
// @Description Create a comment
// @Tags comment
// @Accept json
// @Produce json
// @Param post body models.CreateCommentRequest true "Post"
// @Success 201 {object} models.Comment
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) CreateComment(ctx *gin.Context) {
	var (
		req models.CreateCommentRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	comment, err := h.grpcClient.CommentService().Create(context.Background(), &pb.Comment{
		Description: req.Description,
		UserId:      payload.UserID,
		PostId:      req.PostID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, parseCommentModel(comment))
}

// @Security ApiKeyAuth
// @Router /comments/{id} [put]
// @Summary Update comment
// @Description Update comment
// @Tags comment
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param comment body models.UpdateCommentRequest true "Comment"
// @Success 201 {object} models.UpdateComment
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) UpdateComment(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	var (
		req models.UpdateCommentRequest
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	comment, err := h.grpcClient.CommentService().Update(context.Background(), &pb.Comment{
		Id:          id,
		Description: req.Description,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, parseCommentModel(comment))
}

// @Security ApiKeyAuth
// @Router /comments/{id} [delete]
// @Summary Delete a comment
// @Description Delete a comment
// @Tags comment
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 201 {object} models.ResponseSuccess
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) DeleteComment(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	_, err = h.grpcClient.CommentService().Delete(context.Background(), &pb.DeleteCommentRequest{Id: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Success: "Successfully deleted!",
	})
}

// @Router /comments [get]
// @Summary Get comments
// @Description Get comments
// @Tags comment
// @Accept json
// @Produce json
// @Param filter query models.GetAllCommentsParams false "Filter"
// @Success 201 {object} models.GetAllCommentsResponse
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetAllComments(c *gin.Context) {
	params, err := validateGetAllCommentsParams(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	result, err := h.grpcClient.CommentService().GetAll(context.Background(), &pb.GetAllCommentsParamsReq{
		Limit:  params.Limit,
		Page:   params.Page,
		PostId: params.PostID,
		SortBy: params.SortBy,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, getCommentsResponse(result))
}

func getCommentsResponse(data *pb.GetAllCommentsResponse) *models.GetAllCommentsResponse {
	response := models.GetAllCommentsResponse{
		Comments: make([]*models.Comment, 0),
		Count:    data.Count,
	}

	for _, comment := range data.Comments {
		u := parseCommentModel(comment)
		response.Comments = append(response.Comments, &u)
	}

	return &response
}

func parseCommentModel(comment *pb.Comment) models.Comment {
	return models.Comment{
		ID:          comment.Id,
		Description: comment.Description,
		UserID:      comment.UserId,
		PostID:      comment.PostId,
		CreatedAt:   comment.CreatedAt,
		UpdatedAt:   comment.UpdatedAt,
		User: models.CommentUser{
			FirstName:       comment.User.FirstName,
			Lastname:        comment.User.LastName,
			Email:           comment.User.Email,
			ProfileImageUrl: comment.User.ProfileImageUrl,
		},
	}
}
