package handlers

import (
	"github.com/Poted/raitometer/backend/core-api/internal/database"
	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	db              *sqlx.DB
	projectStore    database.ProjectStore
	userStore       database.UserStore
	calculatorStore database.CalculatorStore
}

func New(
	db *sqlx.DB,
	ps database.ProjectStore,
	us database.UserStore,
	cs database.CalculatorStore,
) *Handlers {
	return &Handlers{
		db:              db,
		projectStore:    ps,
		userStore:       us,
		calculatorStore: cs,
	}
}
