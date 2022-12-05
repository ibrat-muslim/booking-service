package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ibrat-muslim/booking-service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type commentRepo struct {
	db *sqlx.DB
}

func NewComment(db *sqlx.DB) repo.CommentStorageI {
	return &commentRepo{
		db: db,
	}
}

func (cmr *commentRepo) Create(comment *repo.Comment) (*repo.Comment, error) {
	query := `
		INSERT INTO comments (
			post_id,
			user_id,
			description
		) VALUES($1, $2, $3)
		RETURNING id, created_at
	`

	row := cmr.db.QueryRow(
		query,
		comment.PostID,
		comment.UserID,
		comment.Description,
	)

	err := row.Scan(
		&comment.ID,
		&comment.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (cmr *commentRepo) GetAll(params *repo.GetCommentsParams) (*repo.GetCommentsResult, error) {
	result := repo.GetCommentsResult{
		Comments: make([]*repo.Comment, 0),
		Count:    0,
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := " WHERE true "

	if params.PostID != 0 {
		filter += fmt.Sprintf(" AND c.post_id = %d ", params.PostID)
	}

	if params.UserID != 0 {
		filter += fmt.Sprintf(" AND c.user_id = %d ", params.UserID)
	}

	query := `
		SELECT
			c.id,
			c.post_id,
			c.user_id,
			c.description,
			c.created_at,
			c.updated_at,
			u.first_name,
			u.last_name,
			u.email,
			u.profile_image_url
		FROM comments c
		INNER JOIN users u ON u.id = c.user_id
		` + filter + `
		ORDER BY c.created_at DESC
		` + limit

	rows, err := cmr.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var comment repo.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Description,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.User.FirstName,
			&comment.User.LastName,
			&comment.User.Email,
			&comment.User.ProfileImageUrl,
		)
		if err != nil {
			return nil, err
		}

		result.Comments = append(result.Comments, &comment)
	}

	queryCount := `
		SELECT count(1) FROM comments c
		INNER JOIN users u ON u.id = c.user_id ` + filter

	err = cmr.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (cmr *commentRepo) Update(comment *repo.Comment) error {
	query := `
		UPDATE comments SET
			description = $1,
			updated_at = $2
		WHERE id = $3
	`

	result, err := cmr.db.Exec(
		query,
		comment.Description,
		comment.UpdatedAt,
		comment.ID,
	)

	if err != nil {
		return err
	}

	rowsCount, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (cmr *commentRepo) Delete(id int64) error {
	query := `DELETE FROM comments WHERE id = $1`

	result, err := cmr.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsCount, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}
