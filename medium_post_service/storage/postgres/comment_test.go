package postgres_test

import (
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
	"gitlab.com/medium-project/medium_post_service/storage/repo"
)

func createComment(t *testing.T) *repo.Comment {
	post := createPost(t)
	comment, err := dbManager.Comment().Create(&repo.Comment{
		PostID:      post.ID,
		UserID:      1,
		Description: faker.Sentence(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, comment)
	deletePost(t, post.ID)
	return comment
}

func deleteComment(t *testing.T, id int64) {
	err := dbManager.Comment().Delete(id)
	require.NoError(t, err)
}

func TestCreateComment(t *testing.T) {
	c := createComment(t)
	require.NotEmpty(t, c)
	deleteComment(t, c.ID)
}

func TestGetAllComment(t *testing.T) {
	c := createComment(t)
	require.NotEmpty(t, c)
	co, err := dbManager.Comment().GetAll(&repo.GetCommentsParams{
		Limit: 10,
		Page:  1,
	})
	require.GreaterOrEqual(t, len(co.Comments), 1)
	require.NoError(t, err)
	deleteComment(t, c.ID)
}

func TestUpdateComment(t *testing.T) {
	c := createComment(t)
	post := createPost(t)
	require.NotEmpty(t, c)
	co, err := dbManager.Comment().Update(&repo.Comment{
		PostID:      post.ID,
		UserID:      1,
		Description: faker.Sentence(),
	})
	require.NotEmpty(t, co)
	require.NoError(t, err)
	deleteComment(t, c.ID)
	deletePost(t, post.ID)
}

func TestDeleteComment(t *testing.T) {
	c := createComment(t)
	require.NotEmpty(t, c)
	deleteComment(t, c.ID)
}
