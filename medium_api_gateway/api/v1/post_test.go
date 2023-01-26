package v1_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/bxcodec/faker/v4"
// 	"github.com/stretchr/testify/assert"
// 	"gitlab.com/medium-project/medium_api_gateway/api/models"
// )

// func createPost(t *testing.T) *models.Post {
// 	userLogin := loginUser(t)
// 	category := createCategory(t)
// 	resp := httptest.NewRecorder()
// 	arg := models.CreatePostRequest{
// 		Title:       faker.Sentence(),
// 		Description: faker.Sentence(),
// 		CategoryID:  category.ID,
// 	}
// 	payload, err := json.Marshal(arg)
// 	assert.NoError(t, err)
// 	req, err := http.NewRequest("POST", "/v1/posts", bytes.NewBuffer(payload))
// 	assert.NoError(t, err)
// 	req.Header.Add("Authorization", userLogin.AccessToken)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, 200, resp.Code)

// 	body, err := io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var response models.Post
// 	err = json.Unmarshal(body, &response)
// 	assert.NoError(t, err)

// 	return &response
// }

// func deletePost(t *testing.T, postId int64) {
// 	userLogin := loginUser(t)
// 	// * testing delete here
// 	resp := httptest.NewRecorder()
// 	url := fmt.Sprintf("/v1/posts/%d", postId)
// 	req, err := http.NewRequest("DELETE", url, nil)

// 	assert.NoError(t, err)
// 	req.Header.Add("Authorization", userLogin.AccessToken)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, http.StatusOK, resp.Code)

// 	body, err := io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var response models.ResponseSuccess
// 	err = json.Unmarshal(body, &response)
// 	assert.NoError(t, err)

// 	// * testing that post was deleted
// 	resp = httptest.NewRecorder()

// 	url = fmt.Sprintf("/v1/posts/%d", postId)
// 	req, err = http.NewRequest("GET", url, nil)
// 	assert.NoError(t, err)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, http.StatusInternalServerError, resp.Code)

// 	body, err = io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var res models.ResponseError
// 	err = json.Unmarshal(body, &res)
// 	assert.NoError(t, err)
// 	assert.Equal(t, res.Error, "rpc error: code = Internal desc = internal server error: sql: no rows in result set")
// }

// func TestCreatePost(t *testing.T) {
// 	post := createPost(t)
// 	deletePost(t, post.ID)
// }

// func TestGetPost(t *testing.T) {
// 	post := createPost(t)
// 	userLogin := loginUser(t)
// 	resp := httptest.NewRecorder()
// 	url := fmt.Sprintf("/v1/posts/%d", post.ID)
// 	req, err := http.NewRequest("GET", url, nil)
// 	assert.NoError(t, err)

// 	req.Header.Add("Authorization", userLogin.AccessToken)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, 200, resp.Code)

// 	body, err := io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var response models.Post
// 	err = json.Unmarshal(body, &response)
// 	assert.NoError(t, err)
// 	deletePost(t, post.ID)
// }

// func TestUpdatePost(t *testing.T) {
// 	post := createPost(t)
// 	userLogin := loginUser(t)
// 	resp := httptest.NewRecorder()
// 	arg := models.UpdatePostRequest{
// 		Title:       faker.Sentence(),
// 		Description: faker.Sentence(),
// 		ImageUrl:    faker.URL(),
// 		CategoryID:  post.CategoryID,
// 	}
// 	payload, err := json.Marshal(arg)
// 	assert.NoError(t, err)
// 	url := fmt.Sprintf("/v1/posts/%d", post.ID)
// 	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
// 	assert.NoError(t, err)
// 	req.Header.Add("Authorization", userLogin.AccessToken)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, http.StatusOK, resp.Code)

// 	body, err := io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var response models.Post
// 	err = json.Unmarshal(body, &response)
// 	assert.NoError(t, err)

// 	assert.NotZero(t, response.ID)
// 	assert.Equal(t, post.UserID, response.UserID)
// 	assert.Equal(t, arg.Title, response.Title)
// 	assert.Equal(t, arg.Description, response.Description)
// 	assert.Equal(t, arg.CategoryID, response.CategoryID)
// 	assert.Equal(t, arg.ImageUrl, response.ImageUrl)
// 	assert.NotZero(t, response.UpdatedAt)
// 	assert.Equal(t, post.CreatedAt, response.CreatedAt)

// 	deletePost(t, post.ID)
// }

// func TestDeletePost(t *testing.T) {
// 	post := createPost(t)
// 	deletePost(t, post.ID)
// }

// func TestGetAllPosts(t *testing.T) {
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