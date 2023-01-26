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
	pbp "gitlab.com/medium-project/medium_api_gateway/genproto/post_service"
	"gitlab.com/medium-project/medium_api_gateway/pkg/grpc_client/mock_grpc"
)

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arg := models.Category{
		ID:        1,
		Title:     faker.Word(),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	post_service := mock_grpc.NewMockCategoryServiceClient(ctrl)
	post_service.EXPECT().Create(context.Background(), &pbp.Category{
		Title: arg.Title,
	}).Times(1).Return(&pbp.Category{
		Id:        arg.ID,
		Title:     arg.Title,
		CreatedAt: arg.CreatedAt,
	}, nil)

	grpcConn.SetCategoryService(post_service)

	accessToken := mockAuthMiddleware(t, ctrl, "categories", "create")

	payload, err := json.Marshal(arg)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/v1/categories", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Add("Authorization", accessToken)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response models.Category
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, arg.Title, response.Title)
	assert.Equal(t, arg.ID, response.ID)
	assert.Equal(t, arg.CreatedAt, response.CreatedAt)
}

// func deleteCategory(t *testing.T, categoryID int64) {
// 	userLogin := loginSuperAdmin(t)
// 	// * testing delete here
// 	resp := httptest.NewRecorder()
// 	url := fmt.Sprintf("/v1/categories/%d", categoryID)
// 	req, err := http.NewRequest("DELETE", url, nil)

// 	assert.NoError(t, err)
// 	req.Header.Add("Authorization", userLogin.AccessToken)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, 200, resp.Code)

// 	body, err := io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var response models.ResponseSuccess
// 	err = json.Unmarshal(body, &response)
// 	assert.NoError(t, err)
// 	log.Println(response)

// 	// * testing that category was deleted
// 	resp = httptest.NewRecorder()

// 	url = fmt.Sprintf("/v1/categories/%d", categoryID)
// 	req, err = http.NewRequest("GET", url, nil)
// 	assert.NoError(t, err)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, 500, resp.Code)

// 	body, err = io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var res models.ResponseError
// 	err = json.Unmarshal(body, &res)
// 	assert.NoError(t, err)
// 	assert.Equal(t, res.Error, "rpc error: code = Internal desc = internal server error: sql: no rows in result set")
// }

// func TestGetCategory(t *testing.T) {
// 	category := createCategory(t)

// 	resp := httptest.NewRecorder()

// 	url := fmt.Sprintf("/v1/categories/%d", category.ID)
// 	req, err := http.NewRequest("GET", url, nil)
// 	assert.NoError(t, err)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, 200, resp.Code)

// 	body, err := io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var response models.Category
// 	err = json.Unmarshal(body, &response)
// 	assert.NoError(t, err)

// 	assert.Equal(t, category.Title, response.Title)
// 	assert.Equal(t, category.ID, response.ID)
// 	assert.Equal(t, category.CreatedAt, response.CreatedAt)
// 	deleteCategory(t, category.ID)
// }

// func TestUpdateCategory(t *testing.T) {
// 	userLogin := loginSuperAdmin(t)
// 	category := createCategory(t)

// 	resp := httptest.NewRecorder()
// 	arg := models.CreateCategoryRequest{
// 		Title: faker.Sentence(),
// 	}
// 	payload, err := json.Marshal(arg)
// 	assert.NoError(t, err)

// 	url := fmt.Sprintf("/v1/categories/%d", category.ID)
// 	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
// 	assert.NoError(t, err)
// 	req.Header.Add("Authorization", userLogin.AccessToken)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, 200, resp.Code)

// 	body, err := io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var response models.Category
// 	err = json.Unmarshal(body, &response)
// 	assert.NoError(t, err)

// 	assert.Equal(t, category.ID, response.ID)
// 	assert.Equal(t, arg.Title, response.Title)
// 	assert.NotZero(t, category.CreatedAt, response.CreatedAt)
// 	deleteCategory(t, response.ID)
// }

// func TestDeleteCategory(t *testing.T) {
// 	category := createCategory(t)
// 	deleteCategory(t, category.ID)
// }

// func TestGetAllCategories(t *testing.T) {
// 	testCases := []struct {
// 		name          string
// 		query         string
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:  "success code",
// 			query: "?limit=10&page=1",
// 			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusOK, response.Code)
// 			},
// 		},
// 		{
// 			name:  "alpha limit",
// 			query: "?limit=red&page=1",
// 			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusBadRequest, response.Code)
// 			},
// 		},
// 		{
// 			name:  "alpha page",
// 			query: "?limit=10&page=red",
// 			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusBadRequest, response.Code)
// 			},
// 		},
// 		{
// 			name:  "negative limit",
// 			query: "?limit=-10&page=1",
// 			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusInternalServerError, response.Code)
// 			},
// 		},
// 		{
// 			name:  "negative page",
// 			query: "?limit=10&page=-1",
// 			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusInternalServerError, response.Code)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			resp := httptest.NewRecorder()
// 			url := fmt.Sprintf("/v1/categories%s", tc.query)

// 			req, err := http.NewRequest("GET", url, nil)
// 			assert.NoError(t, err)
// 			router.ServeHTTP(resp, req)

// 			tc.checkResponse(t, resp)
// 		})
// 	}
// }
