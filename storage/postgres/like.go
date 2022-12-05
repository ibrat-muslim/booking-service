package postgres

import (
	"database/sql"
	"errors"

	"github.com/ibrat-muslim/booking-service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type likeRepo struct {
	db *sqlx.DB
}

func NewLike(db *sqlx.DB) repo.LikeStorageI {
	return &likeRepo{
		db: db,
	}
}

func (lr *likeRepo) CreateOrUpdate(like *repo.Like) error {
	l, err := lr.Get(like.PostID, like.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		query := `
			INSERT INTO likes (
				post_id,
				user_id,
				status
			) VALUES($1, $2, $3)
			RETURNING id
			`

		_, err := lr.db.Exec(
			query,
			like.PostID,
			like.UserID,
			like.Status,
		)
		if err != nil {
			return err
		}

	} else if l != nil && err == nil {
		if l.Status == like.Status {
			query := `DELETE FROM likes WHERE id = $1`

			_, err := lr.db.Exec(query, l.ID)
			if err != nil {
				return err
			}

		} else {
			query := `UPDATE likes SET status = $1 WHERE id = $2`

			_, err := lr.db.Exec(query, like.Status, l.ID)
			if err != nil {
				return err
			}
		}
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return nil
}

func (l *likeRepo) Get(postID, userID int64) (*repo.Like, error) {
	query := `
		SELECT
			id,
			post_id,
			user_id,
			status
		FROM likes
		WHERE post_id = $1 AND user_id = $2
	`

	var result repo.Like

	err := l.db.Get(&result, query, postID, userID)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (l *likeRepo) GetLikesDislikesCount(postID int64) (*repo.LikesDislikesCountsResult, error) {
	var result repo.LikesDislikesCountsResult

	query := `
		SELECT
			COUNT(1) FILTER (WHERE status = true) as likes_count,
			COUNT(1) FILTER (WHERE status = false) as dislikes_count 
		FROM likes
		WHERE post_id = $1
		`

	err := l.db.Get(&result, query, postID)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
