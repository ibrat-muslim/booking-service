package models

import "time"

type Comment struct {
	ID          int64        `json:"id"`
	PostID      int64        `json:"post_id"`
	UserID      int64        `json:"user_id"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   *time.Time   `json:"updated_at"`
	User        *CommentUser `json:"user"`
}

type CommentUser struct {
	ID              int64   `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	ProfileImageUrl *string `json:"profile_image_url"`
}

type CreateCommentRequest struct {
	PostID      int64  `json:"post_id" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type GetCommentsParams struct {
	Limit  int32 `json:"limit" binding:"required" default:"10"`
	Page   int32 `json:"page" binding:"required" default:"1"`
	PostID int64 `json:"post_id"`
	UserID int64 `json:"user_id"`
}

type GetCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	Count    int32      `json:"count"`
}
