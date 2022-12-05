package repo

type Like struct {
	ID     int64 `db:"id"`
	PostID int64 `db:"post_id"`
	UserID int64 `db:"user_id"`
	Status bool  `db:"status"`
}

type LikesDislikesCountsResult struct {
	LikesCount    int64 `db:"likes_count"`
	DislikesCount int64	`db:"dislikes_count"`
}

type LikeStorageI interface {
	CreateOrUpdate(like *Like) error
	Get(postID, userID int64) (*Like, error)
	GetLikesDislikesCount(postID int64) (*LikesDislikesCountsResult, error)
}
