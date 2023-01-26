package repo

type Like struct {
	ID     int64
	PostID int64
	UserID int64
	Status bool
}

type LikeStorageI interface {
	CreateOrUpdate(like *Like) (*Like, error)
	Get(userID, postID int64) (*Like, error)
	GetLikesDislikesCount(postID int64) (*LikesDislikesCountResult, error) 
}

type LikesDislikesCountResult struct {
	Likes    int64
	Dislikes int64
}