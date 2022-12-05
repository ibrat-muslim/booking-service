package postgres_test

import (
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/ibrat-muslim/booking-service/storage/repo"
	"github.com/stretchr/testify/require"
)

func createCategory(t *testing.T) *repo.Category {
	category, err := strg.Category().Create(&repo.Category{
		Title: faker.Sentence(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, category)

	return category
}

func deleteCategory(id int64, t *testing.T) {
	err := strg.Category().Delete(id)
	require.NoError(t, err)
}

func TestCreateCategory(t *testing.T) {
	c := createCategory(t)
	deleteCategory(c.ID, t)
}

func TestGetCategory(t *testing.T) {
	c := createCategory(t)

	category, err := strg.Category().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	deleteCategory(category.ID, t)
}

func TestGetAllCategories(t *testing.T) {
	c := createCategory(t)

	categories, err := strg.Category().GetAll(&repo.GetCategoriesParams{
		Limit:  10,
		Page:   1,
		Search: c.Title,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(categories.Categories), 1)
	require.GreaterOrEqual(t, int(categories.Count), 1)

	deleteCategory(c.ID, t)
}

func TestUpdateCategory(t *testing.T) {
	c := createCategory(t)

	c.Title = faker.Sentence()

	err := strg.Category().Update(c)
	require.NoError(t, err)

	deleteCategory(c.ID, t)
}

func TestDeleteCategory(t *testing.T) {
	c := createCategory(t)
	deleteCategory(c.ID, t)
}
