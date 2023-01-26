package storage

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/medium-project/medium_post_service/storage/postgres"
	"gitlab.com/medium-project/medium_post_service/storage/repo"
)

type StorageI interface {
	Category() repo.CategoryStorageI
	Post() repo.PostStorageI
	Comment() repo.CommentStorageI
	Like() repo.LikeStorageI
}

type StoragePg struct {
	categoryRepo repo.CategoryStorageI
	postRepo     repo.PostStorageI
	commentRepo  repo.CommentStorageI
	likeRepo     repo.LikeStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &StoragePg{
		categoryRepo: postgres.NewCategory(db),
		postRepo:     postgres.NewPost(db),
		commentRepo:  postgres.NewComment(db),
		likeRepo:     postgres.NewLike(db),
	}
}

func (s *StoragePg) Category() repo.CategoryStorageI {
	return s.categoryRepo
}

func (s *StoragePg) Post() repo.PostStorageI {
	return s.postRepo
}

func (s *StoragePg) Comment() repo.CommentStorageI {
	return s.commentRepo
}

func (s *StoragePg) Like() repo.LikeStorageI {
	return s.likeRepo
}
