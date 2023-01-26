package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arg := models.CreateUserRequest{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Gender:    "male",
		Type:      "superadmin",
		Password:  "ft%2yGYG27t8",
	}

	userService := mock_grpc.NewMockUserServiceClient(ctrl)
	userService.EXPECT().Create(context.Background(), &pbu.User{
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		Password:  arg.Password,
		Type:      arg.Type,
		Gender:    arg.Gender,
	}).Times(1).Return(&pbu.User{
		Id:        1,
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		Type:      arg.Type,
		Gender:    arg.Gender,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil)

	payload, err := json.Marshal(arg)
	assert.NoError(t, err)

	grpcConn.SetUserService(userService)

	accessToken := mockAuthMiddleware(t, ctrl, "users", "create")

	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(payload))
	req.Header.Add("Authorization", accessToken)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response models.User
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, arg.FirstName, response.FirstName)
	assert.Equal(t, arg.Email, response.Email)
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	arg := pbu.User{
		Id:        1,
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Type:      "users",
		Gender:    "male",
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	userService := mock_grpc.NewMockUserServiceClient(ctrl)
	userService.EXPECT().Get(context.Background(), &pbu.IdRequest{
		Id: arg.Id,
	}).Times(1).Return(&arg, nil)

	grpcConn.SetUserService(userService)

	req, _ := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response models.User
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, arg.Id, response.ID)
	assert.Equal(t, arg.Email, response.Email)
}

func TestGetUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arg := pbu.User{
		Id:        1,
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Type:      "users",
		Gender:    "male",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	userService := mock_grpc.NewMockUserServiceClient(ctrl)
	userService.EXPECT().Get(context.Background(), &pbu.IdRequest{
		Id: arg.Id,
	}).Times(1).Return(&arg, nil)

	grpcConn.SetUserService(userService)

	accessToken := mockAuthMiddleware(t, ctrl, "users", "get")

	req, _ := http.NewRequest("GET", "/v1/users/me", nil)
	req.Header.Add("Authorization", accessToken)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response models.User
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arg := models.UpdateUserRequest{
		FirstName:   faker.FirstName(),
		LastName:    faker.LastName(),
		PhoneNumber: faker.Phonenumber(),
		Gender:      "female",
		UserName:    faker.Username(),
	}

	userService := mock_grpc.NewMockUserServiceClient(ctrl)
	userService.EXPECT().Update(context.Background(), &pbu.User{
		Id:          1,
		FirstName:   arg.FirstName,
		LastName:    arg.LastName,
		PhoneNumber: arg.PhoneNumber,
		Gender:      arg.Gender,
		Username:    arg.UserName,
	}).Times(1).Return(&pbu.User{
		Id:          1,
		FirstName:   arg.FirstName,
		LastName:    arg.LastName,
		PhoneNumber: arg.PhoneNumber,
		Gender:      arg.Gender,
		Username:    arg.UserName,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}, nil)

	grpcConn.SetUserService(userService)

	accessToken := mockAuthMiddleware(t, ctrl, "users", "update")

	url := fmt.Sprintf("/v1/users/%d", 1)
	payload, err := json.Marshal(arg)
	assert.NoError(t, err)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", accessToken)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response models.User
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, arg.FirstName, response.FirstName)
	assert.NotZero(t, response.CreatedAt)
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userService := mock_grpc.NewMockUserServiceClient(ctrl)
	userService.EXPECT().Delete(context.Background(), &pbu.IdRequest{
		Id: 1,
	}).Times(1).Return(&emptypb.Empty{}, nil)

	grpcConn.SetUserService(userService)

	accessToken := mockAuthMiddleware(t, ctrl, "users", "delete")
	url := fmt.Sprintf("/v1/users/%d", 1)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", accessToken)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response models.ResponseSuccess
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, response.Success, "Successfully deleted!")
}

func TestGetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arg := pbu.GetAllUsersResponse{
		Count: 2,
		Users: []*pbu.User{
			{
				Id:        1,
				FirstName: faker.FirstName(),
				LastName:  faker.LastName(),
				Email:     faker.Email(),
				Type:      "users",
				Gender:    "male",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			{
				Id:        2,
				FirstName: faker.FirstName(),
				LastName:  faker.LastName(),
				Email:     faker.Email(),
				Type:      "users",
				Gender:    "male",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	}

	userService := mock_grpc.NewMockUserServiceClient(ctrl)
	userService.EXPECT().GetAll(context.Background(), &pbu.GetAllUsersRequest{
		Limit: 10,
		Page:  1,
	}).Times(1).Return(&arg, nil)

	grpcConn.SetUserService(userService)

	testCases := []struct {
		name          string
		query         string
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "success case",
			query: "?limit=10&page=1",
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		{
			name:  "incorrect limit param",
			query: "?limit=sds&page=1",
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/v1/users%s", tc.query)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			tc.checkResponse(t, resp)
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var response models.GetAllUsersResponse
			err = json.Unmarshal(body, &response)
			assert.NoError(t, err)
		})
	}
}
