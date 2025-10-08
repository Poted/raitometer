package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Calculator struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	ProjectID uuid.UUID       `json:"projectID" db:"project_id"`
	Title     string          `json:"title" db:"title"`
	CreatedAt time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time       `json:"updatedAt" db:"updated_at"`
	Modules   json.RawMessage `json:"modules" db:"modules"`
}

type CalculatorModule struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	CalculatorID uuid.UUID       `json:"calculatorID" db:"calculator_id"`
	Title        string          `json:"title" db:"title"`
	DisplayOrder int             `json:"displayOrder" db:"display_order"`
	CreatedAt    time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time       `json:"updatedAt" db:"updated_at"`
	Items        json.RawMessage `json:"items" db:"items"`
}

type CalculatorItem struct {
	ID           uuid.UUID `json:"id" db:"id"`
	ModuleID     uuid.UUID `json:"moduleID" db:"module_id"`
	Description  string    `json:"description" db:"description"`
	PriceType    string    `json:"priceType" db:"price_type"`
	AITag        *string   `json:"aiTag" db:"ai_tag"`
	UnitPrice    float64   `json:"unitPrice" db:"unit_price"`
	Unit         *string   `json:"unit" db:"unit"`
	Quantity     float64   `json:"quantity" db:"quantity"`
	DisplayOrder int       `json:"displayOrder" db:"display_order"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}
