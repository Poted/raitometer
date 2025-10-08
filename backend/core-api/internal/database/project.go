package database

import (
	"database/sql"

	"github.com/Poted/raitometer/backend/core-api/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProjectStore interface {
	Create(project *models.Project) error
	GetByID(id, userID uuid.UUID) (*models.Project, error)
	GetAll(userID uuid.UUID) ([]*models.Project, error)
	Update(project *models.Project) error
	Delete(id, userID uuid.UUID) error
}

type PostgresProjectStore struct {
	DB *sqlx.DB
}

func NewPostgresProjectStore(db *sqlx.DB) *PostgresProjectStore {
	return &PostgresProjectStore{
		DB: db,
	}
}

func (s *PostgresProjectStore) Create(p *models.Project) error {
	query := `INSERT INTO projects (name, description, user_id)
			  VALUES ($1, $2, $3)
			  RETURNING id, created_at, updated_at`

	return s.DB.QueryRowx(query, p.Name, p.Description, p.UserID).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (s *PostgresProjectStore) GetByID(id, userID uuid.UUID) (*models.Project, error) {
	var p models.Project
	query := `SELECT id, user_id, name, description, created_at, updated_at
			  FROM projects
			  WHERE id = $1 AND user_id = $2`

	err := s.DB.Get(&p, query, id, userID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PostgresProjectStore) GetAll(userID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	query := `SELECT id, user_id, name, description, created_at, updated_at 
			  FROM projects 
			  WHERE user_id = $1 
			  ORDER BY created_at DESC`

	err := s.DB.Select(&projects, query, userID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *PostgresProjectStore) Update(p *models.Project) error {
	query := `UPDATE projects
			  SET name = $1, description = $2, updated_at = NOW()
			  WHERE id = $3 AND user_id = $4
			  RETURNING updated_at`

	return s.DB.QueryRowx(query, p.Name, p.Description, p.ID, p.UserID).Scan(&p.UpdatedAt)
}

func (s *PostgresProjectStore) Delete(id, userID uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1 AND user_id = $2`

	result, err := s.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
