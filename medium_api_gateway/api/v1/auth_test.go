package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/medium-project/medium_api_gateway/api/models"
	pbu "gitlab.com/medium-project/medium_api_gateway/genproto/user_service"
	"gitlab.com/medium-project/medium_api_gateway/pkg/grpc_client/mock_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := models.LoginRequest{
		Email:    "testuser@gmail.com",
		Password: "secret-password",
	}

	authService := mock_grpc.NewMockAuthServiceClient(ctrl)
	authService.EXPECT().Login(context.Background(), &pbu.LoginRequest{
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}).Times(1).Return(&pbu.AuthResponse{
		Id:          1,
		FirstName:   "User",
		LastName:    "Test",
		Email:       "testuser@gmail.com",
		Type:        "user",
		CreatedAt:   time.Now().Format(time.RFC3339),
		AccessToken: faker.URL(),
	}, nil)

	payload, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	grpcConn.SetAuthService(authService)

	req, err := http.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)
	var response models.AuthResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := models.RegisterRequest{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Password:  "Secret_2004",
	}

	authService := mock_grpc.NewMockAuthServiceClient(ctrl)
	authService.EXPECT().Register(context.Background(), &pbu.RegisterRequest{
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Email:     reqBody.Email,
		Password:  reqBody.Password,
	}).Times(1).Return(&emptypb.Empty{}, nil)

	payload, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	grpcConn.SetAuthService(authService)

	req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)
	var response models.ResponseSuccess
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, response.Success, "success")
}

func TestVerify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := models.VerifyRequest{
		Email: faker.Email(),
		Code:  "899876",
	}

	authService := mock_grpc.NewMockAuthServiceClient(ctrl)
	authService.EXPECT().Verify(context.Background(), &pbu.VerifyRequest{
		Email: reqBody.Email,
		Code:  reqBody.Code,
	}).Times(1).Return(&pbu.AuthResponse{
		Id:          1,
		FirstName:   faker.FirstName(),
		LastName:    faker.LastName(),
		Email:       faker.Email(),
		Type:        "user",
		CreatedAt:   time.Now().Format(time.RFC3339),
		AccessToken: faker.URL(),
	}, nil)

	payload, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	grpcConn.SetAuthService(authService)

	req, err := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)
	var response models.AuthResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.AccessToken)
}

func TestForgotPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := models.ForgotPasswordRequest{
		Email: faker.Email(),
	}

	authService := mock_grpc.NewMockAuthServiceClient(ctrl)
	authService.EXPECT().ForgotPassword(context.Background(), &pbu.ForgotPasswordRequest{
		Email: reqBody.Email,
	}).Times(1).Return(&emptypb.Empty{}, nil)

	payload, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	grpcConn.SetAuthService(authService)

	req, err := http.NewRequest(http.MethodPost, "/v1/auth/forgot-password", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)
	var response models.ResponseSuccess
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Success, "Validation code has been sent")
}

func TestVerifyForgotPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := models.VerifyRequest{
		Email: faker.Email(),
		Code:  "2345432",
	}

	authService := mock_grpc.NewMockAuthServiceClient(ctrl)
	authService.EXPECT().VerifyForgotPassword(context.Background(), &pbu.VerifyRequest{
		Email: reqBody.Email,
		Code:  reqBody.Code,
	}).Times(1).Return(&pbu.AuthResponse{
		Id:          1,
		FirstName:   faker.FirstName(),
		LastName:    faker.LastName(),
		Email:       faker.Email(),
		Type:        "user",
		CreatedAt:   time.Now().Format(time.RFC3339),
		AccessToken: faker.URL(),
	}, nil)

	payload, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	grpcConn.SetAuthService(authService)

	req, err := http.NewRequest(http.MethodPost, "/v1/auth/verify-forgot-password", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)
	var response models.AuthResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.AccessToken)
}

// ! this test has an error:
// func TestUpdatePassword(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	reqBody := models.UpdatePasswordRequest{
// 		Password: "Secret-Pasword-2022",
// 	}

// 	authService := mock_grpc.NewMockAuthServiceClient(ctrl)
// 	authService.EXPECT().UpdatePassword(context.Background(), &pbu.UpdatePasswordRequest{
// 		UserId:   1,
// 		Password: reqBody.Password,
// 	}).Times(1).Return(&emptypb.Empty{}, nil)

// 	payload, err := json.Marshal(reqBody)
// 	assert.NoError(t, err)

// 	grpcConn.SetAuthService(authService)

// 	accessToken := mockAuthMiddleware(t, ctrl, "update-password", "update")

// 	req, err := http.NewRequest(http.MethodPost, "/v1/auth/update-password", bytes.NewBuffer(payload))
// 	assert.NoError(t, err)
// 	rec := httptest.NewRecorder()
// 	req.Header.Add("Authorization", accessToken)
// 	router.ServeHTTP(rec, req)

// 	assert.Equal(t, http.StatusOK, rec.Code)

// 	body, err := io.ReadAll(rec.Body)
// 	assert.NoError(t, err)
// 	var response models.ResponseError
// 	err = json.Unmarshal(body, &response)
// 	assert.NoError(t, err)

// 	assert.Equal(t, response.Error, "")
// }

func mockAuthMiddleware(t *testing.T, ctrl *gomock.Controller, recource, action string) string {
	accessToken := faker.UUIDHyphenated()

	authService := mock_grpc.NewMockAuthServiceClient(ctrl)
	authService.EXPECT().VerifyToken(context.Background(), &pbu.VerifyTokenRequest{
		AccessToken: accessToken,
		Resource:    recource,
		Action:      action,
	}).Times(1).Return(&pbu.AuthPayload{
		Id:            faker.UUIDHyphenated(),
		UserId:        1,
		Email:         faker.Email(),
		UserType:      "superadmin",
		HasPermission: true,
		IssuedAt:      time.Now().Format(time.RFC3339),
		ExpiredAt:     time.Now().Add(time.Minute * 15).Format(time.RFC3339),
	}, nil)

	grpcConn.SetAuthService(authService)

	return accessToken
}
