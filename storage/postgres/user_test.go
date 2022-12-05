package postgres_test

import (
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/ibrat-muslim/booking-service/storage/repo"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) *repo.User {
	user, err := strg.User().Create(&repo.User{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		DateOfBirth: faker.Date(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		Type:      repo.UserTypeGuest,
	})

	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func deleteUser(id int64, t *testing.T) {
	err := strg.User().Delete(id)
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	u := createUser(t)
	deleteUser(u.ID, t)
}

func TestGetUser(t *testing.T) {
	u := createUser(t)

	user, err := strg.User().Get(u.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	deleteUser(user.ID, t)
}

func TestGetAllUsers(t *testing.T) {
	u := createUser(t)

	users, err := strg.User().GetAll(&repo.GetUsersParams{
		Limit:  10,
		Page:   1,
		Search: u.FirstName,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users.Users), 1)
	require.GreaterOrEqual(t, int(users.Count), 1)

	deleteUser(u.ID, t)
}

func TestUpdateUser(t *testing.T) {
	u := createUser(t)

	u.FirstName = faker.FirstName()
	u.LastName = faker.LastName()
	u.DateOfBirth = faker.Date()
	u.Email = faker.Email()
	u.Password = faker.Password()
	u.Type = repo.UserTypeGuest

	err := strg.User().Update(u)
	require.NoError(t, err)

	deleteUser(u.ID, t)
}

func TestDeleteUser(t *testing.T) {
	u := createUser(t)
	deleteUser(u.ID, t)
}
