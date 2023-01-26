package v1

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/medium-project/medium_api_gateway/api/models"
	"gitlab.com/medium-project/medium_api_gateway/genproto/user_service"
)

// @Router /auth/register [post]
// @Summary Create user with token key and get token key.
// @Description Create user with token key and get token key.
// @Tags register
// @Accept json
// @Produce json
// @Param data body models.RegisterRequest true "Data"
// @Success 200 {object} models.ResponseSuccess
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) Register(ctx *gin.Context) {
	var (
		req models.RegisterRequest
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.logger.WithError(err).Error("failed to bind json in register")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	if !validatePassword(req.Password) {
		ctx.JSON(http.StatusBadRequest, errResponse(ErrWeakPassword))
		return
	}

	user, _ := h.grpcClient.UserService().GetByEmail(context.Background(), &user_service.GetByEmailRequest{
		Email: req.Email,
	})
	if user != nil {
		h.logger.WithError(err).Error("failed to get user by email in register")
		ctx.JSON(http.StatusBadRequest, errResponse(ErrEmailExists))
		return
	}

	_, err = h.grpcClient.AuthService().Register(context.Background(), &user_service.RegisterRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to register user")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Success: "success",
	})
}

func validatePassword(password string) bool {
	var capitalLetter, smallLetter, number, symbol bool
	for i := 0; i < len(password); i++ {
		if password[i] >= 65 && password[i] <= 90 {
			capitalLetter = true
		} else if password[i] >= 97 && password[i] <= 122 {
			smallLetter = true
		} else if password[i] >= 48 && password[i] <= 57 {
			number = true
		} else {
			symbol = true
		}
	}
	return capitalLetter && smallLetter && number && symbol
}

// @Router /auth/verify [post]
// @Summary Create user with token key and get token key.
// @Description Create user with token key and get token key.
// @Tags register
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) Verify(ctx *gin.Context) {
	var (
		req models.VerifyRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("failed to bind JSON in verify")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := h.grpcClient.AuthService().Verify(context.Background(), &user_service.VerifyRequest{
		Email: req.Email,
		Code:  req.Code,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to verify user")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.AuthResponse{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Type:        user.Type,
		CreatedAt:   user.CreatedAt,
		AccessToken: user.AccessToken,
	})
}

// @Router /auth/login [post]
// @Summary Login User
// @Description Login User
// @Tags register
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) Login(ctx *gin.Context) {
	var (
		req models.LoginRequest
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.logger.WithError(err).Error("failed to bind JSON in login")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := h.grpcClient.AuthService().Login(context.Background(), &user_service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to login")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.AuthResponse{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Type:        user.Type,
		CreatedAt:   user.CreatedAt,
		AccessToken: user.AccessToken,
	})

}

// @Router /auth/forgot-password [post]
// @Summary Forgot  password
// @Description Forgot  password
// @Tags forgot_password
// @Accept json
// @Produce json
// @Param data body models.ForgotPasswordRequest true "Data"
// @Success 200 {object} models.ResponseSuccess
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) ForgotPassword(ctx *gin.Context) {
	var (
		req models.ForgotPasswordRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("failed to bind JSON in forgotpassword")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	_, err := h.grpcClient.AuthService().ForgotPassword(context.Background(), &user_service.ForgotPasswordRequest{
		Email: req.Email,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed in forgot password")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Success: "Validation code has been sent",
	})
}

// @Router /auth/verify-forgot-password [post]
// @Summary Verify forgot password
// @Description Verify forgot password
// @Tags forgot_password
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) VerifyForgotPassword(ctx *gin.Context) {
	var (
		req *models.VerifyRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("failed to bind JSON in verifyforgotpassword")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	result, err := h.grpcClient.AuthService().VerifyForgotPassword(context.Background(), &user_service.VerifyRequest{
		Email: req.Email,
		Code:  req.Code,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed in verifyforgotpassword")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.AuthResponse{
		Id:          result.Id,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		Type:        result.Type,
		CreatedAt:   result.CreatedAt,
		AccessToken: result.AccessToken,
	})
}

// @Security ApiKeyAuth
// @Router /auth/update-password [post]
// @Summary Update password
// @Description Update password
// @Tags forgot_password
// @Accept json
// @Produce json
// @Param data body models.UpdatePasswordRequest true "Data"
// @Success 200 {object} models.ResponseSuccess
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) UpdatePassword(ctx *gin.Context) {
	var (
		req models.UpdatePasswordRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("failed to bind JSON in updatepassword")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		h.logger.WithError(err).Error("failed to get payload in update password")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	log.Println(payload)
	_, err = h.grpcClient.AuthService().UpdatePassword(context.Background(), &user_service.UpdatePasswordRequest{
		UserId:   payload.UserID,
		Password: req.Password,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to update password")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Success: "Password has been updated!",
	})
}
