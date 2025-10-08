package database

import (
	"github.com/Poted/raitometer/backend/core-api/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CalculatorStore interface {
	Create(calc *models.Calculator) error
	GetByID(id, userID uuid.UUID) (*models.Calculator, error)
	GetFullByID(id, userID uuid.UUID) (*models.Calculator, error)
	CreateModule(module *models.CalculatorModule) error
	GetModuleByID(id, userID uuid.UUID) (*models.CalculatorModule, error)
	CreateItem(item *models.CalculatorItem) error
}

type PostgresCalculatorStore struct {
	DB *sqlx.DB
}

func NewPostgresCalculatorStore(db *sqlx.DB) *PostgresCalculatorStore {
	return &PostgresCalculatorStore{
		DB: db,
	}
}

func (s *PostgresCalculatorStore) Create(c *models.Calculator) error {
	query := `INSERT INTO calculators (project_id, title)
			  VALUES ($1, $2)
			  RETURNING id, created_at, updated_at`

	return s.DB.QueryRowx(query, c.ProjectID, c.Title).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (s *PostgresCalculatorStore) GetByID(id, userID uuid.UUID) (*models.Calculator, error) {
	var c models.Calculator
	query := `SELECT c.* FROM calculators c
			  JOIN projects p ON c.project_id = p.id
			  WHERE c.id = $1 AND p.user_id = $2`

	err := s.DB.Get(&c, query, id, userID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *PostgresCalculatorStore) GetFullByID(id, userID uuid.UUID) (*models.Calculator, error) {
	var calculator models.Calculator
	query := `
		SELECT
			c.id, c.project_id, c.title, c.created_at, c.updated_at,
			COALESCE(
				(SELECT json_agg(m) FROM (
					SELECT
						m.id, m.calculator_id, m.title, m.display_order, m.created_at, m.updated_at,
						COALESCE(
							(SELECT json_agg(i) FROM (
								SELECT i.id, i.module_id, i.description, i.price_type, i.ai_tag, i.unit_price, i.unit, i.quantity, i.display_order, i.created_at, i.updated_at
								FROM calculator_items i
								WHERE i.module_id = m.id
								ORDER BY i.display_order
							) AS i), '[]'::json
						) as items
					FROM calculator_modules m
					WHERE m.calculator_id = c.id
					ORDER BY m.display_order
				) AS m), '[]'::json
			) as modules
		FROM calculators c
		JOIN projects p ON c.project_id = p.id
		WHERE c.id = $1 AND p.user_id = $2
	`

	err := s.DB.Get(&calculator, query, id, userID)
	if err != nil {
		return nil, err
	}

	return &calculator, nil
}

func (s *PostgresCalculatorStore) CreateModule(m *models.CalculatorModule) error {
	query := `INSERT INTO calculator_modules (calculator_id, title, display_order)
			  VALUES ($1, $2, $3)
			  RETURNING id, created_at, updated_at`

	return s.DB.QueryRowx(query, m.CalculatorID, m.Title, m.DisplayOrder).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (s *PostgresCalculatorStore) GetModuleByID(id, userID uuid.UUID) (*models.CalculatorModule, error) {
	var m models.CalculatorModule
	query := `SELECT m.* FROM calculator_modules m
			  JOIN calculators c ON m.calculator_id = c.id
			  JOIN projects p ON c.project_id = p.id
			  WHERE m.id = $1 AND p.user_id = $2`

	err := s.DB.Get(&m, query, id, userID)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *PostgresCalculatorStore) CreateItem(i *models.CalculatorItem) error {
	query := `INSERT INTO calculator_items (module_id, description, price_type, ai_tag, unit_price, unit, quantity, display_order)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			  RETURNING id, created_at, updated_at`

	return s.DB.QueryRowx(query, i.ModuleID, i.Description, i.PriceType, i.AITag, i.UnitPrice, i.Unit, i.Quantity, i.DisplayOrder).Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
}
