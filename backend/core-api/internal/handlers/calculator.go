package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/Poted/raitometer/backend/core-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (h *Handlers) CreateCalculatorHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	projectIDStr := chi.URLParam(r, "projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Bad Request: invalid project ID", http.StatusBadRequest)
		return
	}

	_, err = h.projectStore.GetByID(projectID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found: project not found", http.StatusNotFound)
		} else {
			log.Printf("error getting project by id: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		Title string `json:"title"`
	}

	err = h.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Bad Request: could not decode JSON", http.StatusBadRequest)
		return
	}

	calc := &models.Calculator{
		ProjectID: projectID,
		Title:     input.Title,
	}

	err = h.calculatorStore.Create(calc)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) && pqError.Code == "23505" {
			http.Error(w, "Conflict: a calculator for this project already exists", http.StatusConflict)
		} else {
			log.Printf("error creating calculator: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	err = h.writeJSON(w, http.StatusCreated, calc, nil)
	if err != nil {
		log.Printf("error writing create calculator response: %v", err)
	}
}

func (h *Handlers) GetFullCalculatorHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	calculatorIDStr := chi.URLParam(r, "calculatorID")
	calculatorID, err := uuid.Parse(calculatorIDStr)
	if err != nil {
		http.Error(w, "Bad Request: invalid calculator ID", http.StatusBadRequest)
		return
	}

	calculator, err := h.calculatorStore.GetFullByID(calculatorID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found: calculator not found", http.StatusNotFound)
		} else {
			log.Printf("error getting full calculator by id: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	err = h.writeJSON(w, http.StatusOK, calculator, nil)
	if err != nil {
		log.Printf("error writing get full calculator response: %v", err)
	}
}

func (h *Handlers) CreateModuleHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	calculatorIDStr := chi.URLParam(r, "calculatorID")
	calculatorID, err := uuid.Parse(calculatorIDStr)
	if err != nil {
		http.Error(w, "Bad Request: invalid calculator ID", http.StatusBadRequest)
		return
	}

	_, err = h.calculatorStore.GetByID(calculatorID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found: calculator not found", http.StatusNotFound)
		} else {
			log.Printf("error getting calculator by id: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		Title        string `json:"title"`
		DisplayOrder int    `json:"displayOrder"`
	}

	err = h.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Bad Request: could not decode JSON", http.StatusBadRequest)
		return
	}

	module := &models.CalculatorModule{
		CalculatorID: calculatorID,
		Title:        input.Title,
		DisplayOrder: input.DisplayOrder,
	}

	err = h.calculatorStore.CreateModule(module)
	if err != nil {
		log.Printf("error creating calculator module: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.writeJSON(w, http.StatusCreated, module, nil)
	if err != nil {
		log.Printf("error writing create module response: %v", err)
	}
}

func (h *Handlers) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	moduleIDStr := chi.URLParam(r, "moduleID")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		http.Error(w, "Bad Request: invalid module ID", http.StatusBadRequest)
		return
	}

	_, err = h.calculatorStore.GetModuleByID(moduleID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found: module not found", http.StatusNotFound)
		} else {
			log.Printf("error getting module by id: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	var input models.CalculatorItem
	err = h.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Bad Request: could not decode JSON", http.StatusBadRequest)
		return
	}

	item := &models.CalculatorItem{
		ModuleID:     moduleID,
		Description:  input.Description,
		PriceType:    input.PriceType,
		AITag:        input.AITag,
		UnitPrice:    input.UnitPrice,
		Unit:         input.Unit,
		Quantity:     input.Quantity,
		DisplayOrder: input.DisplayOrder,
	}

	err = h.calculatorStore.CreateItem(item)
	if err != nil {
		log.Printf("error creating calculator item: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.writeJSON(w, http.StatusCreated, item, nil)
	if err != nil {
		log.Printf("error writing create item response: %v", err)
	}
}
