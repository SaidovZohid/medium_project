package v1

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/medium-project/medium_api_gateway/api/models"
	pbu "gitlab.com/medium-project/medium_api_gateway/genproto/user_service"
)

// @Security ApiKeyAuth
// @Router /users [post]
// @Summary Create a user
// @Description Create a user
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User"
// @Success 201 {object} models.User
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		req models.CreateUserRequest
	)
	err := c.ShouldBindJSON(&req)

	if err != nil {
		h.logger.WithError(err).Error("failed to bind JSON in create user")
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	resp, err := h.grpcClient.UserService().Create(context.Background(), &pbu.User{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Gender:          req.Gender,
		Password:        req.Password,
		Username:        req.UserName,
		ProfileImageUrl: req.ProfileImageUrl,
		Type:            req.Type,
	})

	if err != nil {
		h.logger.WithError(err).Error("failed to create user")
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseToUserModel(resp))
}

// @Router /users/{id} [get]
// @Summary Get user
// @Description Get user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 201 {object} models.User
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
// @Failure 404 {object} models.ResponseError
func (h *handlerV1) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.WithError(err).Error("failed to bind JSON in create user")
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	resp, err := h.grpcClient.UserService().Get(context.Background(), &pbu.IdRequest{
		Id: id,
	})

	if err != nil {
		h.logger.WithError(err).Error("failed to get user")
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseToUserModel(resp))
}

func parseToUserModel(user *pbu.User) models.User {
	return models.User{
		ID:              user.Id,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Gender:          user.Gender,
		UserName:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt,
	}
}

// @Security ApiKeyAuth
// @Router /users/me [get]
// @Summary Get user by token
// @Description Get user by token
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
// @Failure 404 {object} models.ResponseError
func (h *handlerV1) GetUserProfile(c *gin.Context) {
	payload, err := h.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	resp, err := h.grpcClient.UserService().Get(context.Background(), &pbu.IdRequest{
		Id: payload.UserID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseToUserModel(resp))
}

// @Security ApiKeyAuth
// @Router /users/{id} [put]
// @Summary Update user
// @Description Update user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param User body models.UpdateUserRequest true "User"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		req models.UpdateUserRequest
	)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.WithError(err).Error("failed to parse user id in update user")
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	err = c.ShouldBindJSON(&req)

	if err != nil {
		h.logger.WithError(err).Error("failed to bind JSON in update user")
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	result, err := h.grpcClient.UserService().Update(context.Background(), &pbu.User{
		Id:              id,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Gender:          req.Gender,
		Username:        req.UserName,
		ProfileImageUrl: req.ProfileImageUrl,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to update user")
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseToUserModel(result))
}

// @Security ApiKeyAuth
// @Router /users/{id} [delete]
// @Summary Delete user
// @Description Delete user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 201 {object} models.ResponseSuccess
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	_, err = h.grpcClient.UserService().Delete(context.Background(), &pbu.IdRequest{Id: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.ResponseSuccess{
		Success: "Successfully deleted!",
	})
}

// @Router /users [get]
// @Summary Get user by giving limit, page and search for something.
// @Description Get user by giving limit, page and search for something.
// @Tags user
// @Accept json
// @Produce json
// @Param filter query models.GetAllParams false "Filter"
// @Success 201 {object} models.GetAllUsersResponse
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	params, err := validateGetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	result, err := h.grpcClient.UserService().GetAll(context.Background(), &pbu.GetAllUsersRequest{
		Limit:  int32(params.Limit),
		Page:   int32(params.Page),
		Search: params.Search,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, getUsersResponse(result))
}

func getUsersResponse(data *pbu.GetAllUsersResponse) *models.GetAllUsersResponse {
	response := models.GetAllUsersResponse{
		Users: make([]*models.User, 0),
		Count: data.Count,
	}

	for _, user := range data.Users {
		u := parseToUserModel(user)
		response.Users = append(response.Users, &u)
	}

	return &response
}
