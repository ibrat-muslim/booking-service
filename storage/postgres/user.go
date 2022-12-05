package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ibrat-muslim/booking-service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) Create(user *repo.User) (*repo.User, error) {
	query := `
		INSERT INTO users (
			first_name,
			last_name,
			dob,
			email,
			phone_number,
			gender,
			password,
			profile_image_url,
			address,
			type
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`

	row := ur.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.Email,
		user.PhoneNumber,
		user.Gender,
		user.Password,
		user.ProfileImageUrl,
		user.Address,
		user.Type,
	)

	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) Get(id int64) (*repo.User, error) {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			dob,
			email,
			phone_number,
			gender,
			password,
			profile_image_url,
			address,
			type,
			created_at
		FROM users
		WHERE id = $1
	`

	var result repo.User

	err := ur.db.Get(&result, query, id)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) GetByEmail(email string) (*repo.User, error) {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			dob,
			email,
			phone_number,
			gender,
			password,
			profile_image_url,
			address,
			type,
			created_at
		FROM users
		WHERE email = $1
	`

	var result repo.User

	err := ur.db.Get(&result, query, email)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) GetAll(params *repo.GetUsersParams) (*repo.GetUsersResult, error) {
	result := repo.GetUsersResult{
		Users: make([]*repo.User, 0),
		Count: 0,
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := ""

	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
				WHERE first_name ILIKE '%s' OR last_name ILIKE '%s' OR phone_number ILIKE '%s' 
				OR email ILIKE '%s'`,
			str, str, str, str,
		)
	}

	query := `
		SELECT
			id,
			first_name,
			last_name,
			dob,
			email,
			phone_number,
			gender,
			password,
			profile_image_url,
			address,
			type,
			created_at
		FROM users
		` + filter + `
		ORDER BY created_at DESC
		` + limit

	err := ur.db.Select(&result.Users, query)

	if err != nil {
		return nil, err
	}

	queryCount := `SELECT count(1) FROM users ` + filter

	err = ur.db.Get(&result.Count, queryCount)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) Update(user *repo.User) error {
	query := `
		UPDATE users SET
			first_name = $1,
			last_name = $2,
			dob = $3,
			email = $4,
			phone_number = $5,
			gender = $6,
			password = $7,
			profile_image_url = $8,
			address = $9,
			type = $10
		WHERE id = $11
	`

	result, err := ur.db.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.Email,
		user.PhoneNumber,
		user.Gender,
		user.Password,
		user.ProfileImageUrl,
		user.Address,
		user.Type,
		user.ID,
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

func (ur *userRepo) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	resutl, err := ur.db.Exec(query, id)

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

func (ur *userRepo) UpdatePassword(req *repo.UpdatePassword) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`

	_, err := ur.db.Exec(
		query,
		req.Password,
		req.UserID,
	)

	if err != nil {
		return err
	}

	return nil
}
