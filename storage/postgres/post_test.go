package postgres_test

import (
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/ibrat-muslim/booking-service/storage/repo"
	"github.com/stretchr/testify/require"
)

func createPost(t *testing.T) *repo.Post {
	user := createUser(t)
	category := createCategory(t)

	post, err := strg.Post().Create(&repo.Post{
		Title:       faker.Sentence(),
		Description: faker.Sentence(),
		UserID:      user.ID,
		CategoryID:  category.ID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, post)

	return post
}

func deletePost(id int64, t *testing.T) {
	err := strg.Post().Delete(id)
	require.NoError(t, err)
}

func TestCreatePost(t *testing.T) {
	p := createPost(t)
	deletePost(p.ID, t)
}

func TestGetPost(t *testing.T) {
	p := createPost(t)

	post, err := strg.Post().Get(p.ID)
	require.NoError(t, err)
	require.NotEmpty(t, post)

	deletePost(post.ID, t)
}

func TestGetAllPosts(t *testing.T) {
	p := createPost(t)

	posts, err := strg.Post().GetAll(&repo.GetPostsParams{
		Limit:  10,
		Page:   1,
		Search: p.Title,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(posts.Posts), 1)
	require.GreaterOrEqual(t, int(posts.Count), 1)

	deletePost(p.ID, t)
}

func TestUpdatePost(t *testing.T) {
	p := createPost(t)
	category := createCategory(t)

	p.Title = faker.Sentence()
	p.Description = faker.Sentence()
	p.CategoryID = category.ID

	err := strg.Post().Update(p)
	require.NoError(t, err)

	deletePost(p.ID, t)
}

func TestDeletePost(t *testing.T) {
	p := createPost(t)
	deletePost(p.ID, t)
}
