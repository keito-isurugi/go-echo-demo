package repository

import (
	"database/sql"
	"errors"
	"go-echo-demo/internal/domain"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) ValidateCredentials(email, password string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, name, email, password FROM users WHERE email = $1 AND password = $2`

	err := r.db.QueryRow(query, email, password).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	return &user, nil
}
