package store

import (
	"errors"

	"github.com/mcfly722/experiments/go-rest/model"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if len(u.Password) > 0 {
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
		if err != nil {
			return nil, err
		}

		if err := r.store.db.QueryRow(
			"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
			u.Email,
			encryptedPassword,
		).Scan(&u.ID); err != nil {
			return nil, err
		}
		return u, nil
	} else {
		return nil, errors.New("user password could not be empty")
	}
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		return nil, err
	}

	return u, nil
}
