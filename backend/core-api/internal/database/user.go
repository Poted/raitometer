package database

import (
	"github.com/Poted/raitometer/backend/core-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserStore interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

type PostgresUserStore struct {
	DB *sqlx.DB
}

func NewPostgresUserStore(db *sqlx.DB) *PostgresUserStore {
	return &PostgresUserStore{
		DB: db,
	}
}

func (s *PostgresUserStore) Create(u *models.User) error {
	query := `INSERT INTO users (email, password_hash)
			  VALUES ($1, $2)
			  RETURNING id, created_at, updated_at`

	return s.DB.QueryRowx(query, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (s *PostgresUserStore) GetByEmail(email string) (*models.User, error) {
	var u models.User
	query := `SELECT id, email, password_hash, created_at, updated_at
			  FROM users
			  WHERE email = $1`

	err := s.DB.Get(&u, query, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
