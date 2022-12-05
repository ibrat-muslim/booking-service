package repo

import "time"

type Comment struct {
	ID          int64      `db:"id"`
	PostID      int64      `db:"post_id"`
	UserID      int64      `db:"user_id"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	User        struct {
		FirstName       string  `db:"first_name"`
		LastName        string  `db:"last_name"`
		Email           string  `db:"email"`
		ProfileImageUrl *string `db:"profile_image_url"`
	}
}

type GetCommentsParams struct {
	Limit  int32 `db:"limit"`
	Page   int32 `db:"page"`
	PostID int64 `db:"post_id"`
	UserID int64 `db:"user_id"`
}

type GetCommentsResult struct {
	Comments []*Comment `db:"comments"`
	Count    int32      `db:"count"`
}

type CommentStorageI interface {
	Create(comment *Comment) (*Comment, error)
	GetAll(params *GetCommentsParams) (*GetCommentsResult, error)
	Update(comment *Comment) error
	Delete(id int64) error
}
