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
// @Router /categories [post]
// @Summary Create a category
// @Description Create a category
// @Tags category
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Category"
// @Success 200 {object} models.Category
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) CreateCategory(ctx *gin.Context) {
	var (
		req models.CreateCategoryRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	category, err := h.grpcClient.CategoryService().Create(context.Background(), &pb.Category{
		Title: req.Title,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Category{
		ID:        category.Id,
		Title:     category.Title,
		CreatedAt: category.CreatedAt,
	})
}

// @Router /categories/{id} [get]
// @Summary Get category by it's id
// @Description Get category by it's id
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 201 {object} models.Category
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetCategory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	category, err := h.grpcClient.CategoryService().Get(context.Background(), &pb.GetCategoryRequest{Id: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Category{
		ID:        category.Id,
		Title:     category.Title,
		CreatedAt: category.CreatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /categories/{id} [put]
// @Summary Update category by it's id
// @Description Update category by it's id
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param category body models.CreateCategoryRequest true "Category"
// @Success 201 {object} models.Category
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) UpdateCategory(ctx *gin.Context) {
	var (
		req models.CreateCategoryRequest
		err error
	)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	category, err := h.grpcClient.CategoryService().Update(context.Background(), &pb.Category{
		Id:    id,
		Title: req.Title,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Category{
		ID:        category.Id,
		Title:     category.Title,
		CreatedAt: category.CreatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /categories/{id} [delete]
// @Summary Delete category by it's id
// @Description Delete category by it's id
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 201 {object} models.ResponseSuccess
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	_, err = h.grpcClient.CategoryService().Delete(context.Background(), &pb.GetCategoryRequest{Id: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Success: "Succesfully deleted!",
	})
}

// @Router /categories [get]
// @Summary Get categories
// @Description Get category
// @Tags category
// @Accept json
// @Produce json
// @Param filter query models.GetAllParams false "Filter"
// @Success 201 {object} models.GetCategoriesResponse
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetAllCategories(ctx *gin.Context) {
	params, err := validateGetAllParams(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	category, err := h.grpcClient.CategoryService().GetAll(context.Background(), &pb.GetAllCategoryParamsReq{
		Limit:  int32(params.Limit),
		Page:   int32(params.Page),
		Search: params.Search,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getCategoriesResponse(category))
}

func getCategoriesResponse(categories *pb.GetAllCategoryResponse) *models.GetCategoriesResponse {
	response := models.GetCategoriesResponse{
		Categories: make([]*models.Category, 0),
		Count:      int64(categories.Count),
	}

	for _, c := range categories.Categories {
		response.Categories = append(response.Categories, &models.Category{
			ID:        c.Id,
			Title:     c.Title,
			CreatedAt: c.CreatedAt,
		})
	}
	return &response
}
