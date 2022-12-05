package repo

import "time"

type Category struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
}

type GetCategoriesParams struct {
	Limit  int32  `db:"limit"`
	Page   int32  `db:"page"`
	Search string `db:"search"`
}

type GetCategoriesResult struct {
	Categories []*Category `db:"categories"`
	Count      int32       `db:"count"`
}

type CategoryStorageI interface {
	Create(category *Category) (*Category, error)
	Get(id int64) (*Category, error)
	GetAll(params *GetCategoriesParams) (*GetCategoriesResult, error)
	Update(category *Category) error
	Delete(id int64) error
}
