package repo

import "time"

type Post struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	ImageUrl    *string    `db:"image_url"`
	UserID      int64      `db:"user_id"`
	CategoryID  int64      `db:"category_id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	ViewsCount  int32      `db:"views_count"`
	LikeInfo    struct {
		LikesCount    int64 `db:"likes_count"`
		DisLikesCount int64 `db:"dislikes_count"`
	}
}

type GetPostsParams struct {
	Limit      int32  `db:"limit"`
	Page       int32  `db:"page"`
	Search     string `db:"search"`
	UserID     int64  `db:"user_id"`
	CategoryID int64  `db:"category_id"`
	SortByDate string `db:"sort_by_date"`
}

type GetPostsResult struct {
	Posts []*Post `db:"posts"`
	Count int32   `db:"count"`
}

type PostStorageI interface {
	Create(post *Post) (*Post, error)
	Get(id int64) (*Post, error)
	GetAll(params *GetPostsParams) (*GetPostsResult, error)
	Update(post *Post) error
	Delete(id int64) error
}
