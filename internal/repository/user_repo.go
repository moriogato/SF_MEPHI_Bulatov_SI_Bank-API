package repository

import (
	"bank-api/internal/models"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID)
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash FROM users WHERE email=$1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}
