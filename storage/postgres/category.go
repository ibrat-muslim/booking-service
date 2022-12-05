package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ibrat-muslim/booking-service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type categoryRepo struct {
	db *sqlx.DB
}

func NewCategory(db *sqlx.DB) repo.CategoryStorageI {
	return &categoryRepo{
		db: db,
	}
}

func (cr *categoryRepo) Create(category *repo.Category) (*repo.Category, error) {
	query := `
		INSERT INTO categories (
			title
		) VALUES($1)
		RETURNING id, created_at
	`

	row := cr.db.QueryRow(
		query,
		category.Title,
	)

	err := row.Scan(
		&category.ID,
		&category.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (cr *categoryRepo) Get(id int64) (*repo.Category, error) {
	query := `
		SELECT
			id,
			title,
			created_at
		FROM categories
		WHERE id = $1
	`

	var result repo.Category

	err := cr.db.Get(&result, query, id)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (cr *categoryRepo) GetAll(params *repo.GetCategoriesParams) (*repo.GetCategoriesResult, error) {
	result := repo.GetCategoriesResult{
		Categories: make([]*repo.Category, 0),
		Count:      0,
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := ""

	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
				WHERE title ILIKE '%s'`, str,
		)
	}

	query := `
		SELECT
			id,
			title,
			created_at
		FROM categories
		` + filter + `
		ORDER BY created_at DESC
		` + limit

	err := cr.db.Select(&result.Categories, query)

	if err != nil {
		return nil, err
	}

	queryCount := `SELECT count(1) FROM categories ` + filter

	err = cr.db.Get(&result.Count, queryCount)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (cr *categoryRepo) Update(category *repo.Category) error {
	query := `
		UPDATE categories SET
			title = $1
		WHERE id = $2
	`

	result, err := cr.db.Exec(
		query,
		category.Title,
		category.ID,
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

func (cr *categoryRepo) Delete(id int64) error {
	query := `DELETE FROM posts WHERE category_id = $1`

	_, err := cr.db.Exec(query, id)

	if err != nil {
		return err
	}

	query = `DELETE FROM categories WHERE id = $1`

	resutl, err := cr.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsCount, err := resutl.RowsAffected()

	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}
