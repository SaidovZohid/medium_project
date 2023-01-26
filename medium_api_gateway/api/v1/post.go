package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/medium-project/medium_api_gateway/api/models"
	"gitlab.com/medium-project/medium_api_gateway/genproto/post_service"
)

// @Security ApiKeyAuth
// @Router /posts [post]
// @Summary Create a post
// @Description Create a post
// @Tags post
// @Accept json
// @Produce json
// @Param post body models.CreatePostRequest true "Post"
// @Success 201 {object} models.Post
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
// @Failure 401 {object} models.ResponseError
func (h *handlerV1) CreatePost(ctx *gin.Context) {
	var (
		req models.CreatePostRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	post, err := h.grpcClient.PostService().Create(context.Background(), &post_service.Post{
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		UserId:      payload.UserID,
		CategoryId:  req.CategoryID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.Post{
		ID:          post.Id,
		Title:       post.Title,
		Description: post.Description,
		ImageUrl:    post.ImageUrl,
		UserID:      post.UserId,
		CategoryID:  post.CategoryId,
		CreatedAt:   post.CreatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /posts/{id} [get]
// @Summary Get a post with it's id
// @Description Create a post with it's id
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 201 {object} models.Post
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) GetPost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	res, err := h.grpcClient.PostService().Get(context.Background(), &post_service.GetPostRequest{Id: id})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	post := parsePostModel(res)

	likesInfo, err := h.grpcClient.LikeService().GetLikesDislikesCount(context.Background(), &post_service.GetLikesRequest{PostId: post.ID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	post.PostLikeInfo = models.PostLikeInfo{
		LikesCount:    likesInfo.Likes,
		DislikesCount: likesInfo.Dislikes,
	}
	ctx.JSON(http.StatusOK, post)
}

// @Security ApiKeyAuth
// @Router /posts/{id} [put]
// @Summary Update post with it's id as param
// @Description Update post with it's id as param
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param post body models.UpdatePostRequest true "Post"
// @Success 200 {object} models.Post
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
// @Failure 401 {object} models.ResponseError
func (h *handlerV1) UpdatePost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var (
		req models.UpdatePostRequest
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	post, err := h.grpcClient.PostService().Update(context.Background(), &post_service.Post{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		UserId:      payload.UserID,
		CategoryId:  req.CategoryID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.Post{
		ID:          post.Id,
		Title:       post.Title,
		Description: post.Description,
		ImageUrl:    post.ImageUrl,
		UserID:      post.UserId,
		CategoryID:  post.CategoryId,
		ViewsCount:  int32(post.ViewsCount),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /posts/{id} [delete]
// @Summary Delete a post
// @Description Create a post
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 201 {object} models.ResponseSuccess
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) DeletePost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	_, err = h.grpcClient.PostService().Delete(context.Background(), &post_service.GetPostRequest{Id: id})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Success: "Successfully deleted!",
	})
}

// @Router /posts [get]
// @Summary Get posts by giving limit, page and search for something.
// @Description Get posts by giving limit, page and search for something.
// @Tags post
// @Accept json
// @Produce json
// @Param filter query models.GetAllPostsParams false "Filter"
// @Success 200 {object} models.GetAllPostsResponse
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) GetAllPosts(c *gin.Context) {
	params, err := validateGetAllPostsParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	result, err := h.grpcClient.PostService().GetAll(context.Background(), &post_service.GetPostsParamsReq{
		Limit:      params.Limit,
		Page:       params.Page,
		Search:     params.Search,
		UserId:     params.UserID,
		CategoryId: params.CategoryID,
		SortByDate: params.SortByDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, getPostsResponse(result))
}

func getPostsResponse(data *post_service.GetAllPostResponse) *models.GetAllPostsResponse {
	response := models.GetAllPostsResponse{
		Posts: make([]*models.Post, 0),
		Count: data.Count,
	}

	for _, post := range data.Posts {
		u := parsePostModel(post)
		response.Posts = append(response.Posts, &u)
	}

	return &response
}

func parsePostModel(post *post_service.Post) models.Post {
	return models.Post{
		ID:          post.Id,
		Title:       post.Title,
		Description: post.Title,
		ImageUrl:    post.ImageUrl,
		ViewsCount:  int32(post.ViewsCount),
		UserID:      post.UserId,
		CategoryID:  post.CategoryId,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}
