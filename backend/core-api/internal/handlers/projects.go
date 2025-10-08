package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Poted/raitometer/backend/core-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var input struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	err = h.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Bad Request: could not decode JSON", http.StatusBadRequest)
		return
	}

	project := &models.Project{
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
	}

	err = h.projectStore.Create(project)
	if err != nil {
		log.Printf("error creating project: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.writeJSON(w, http.StatusCreated, project, nil)
	if err != nil {
		log.Printf("error writing create project response: %v", err)
	}
}

func (h *Handlers) GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "projectID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Bad Request: invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := h.projectStore.GetByID(id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found: project not found", http.StatusNotFound)
		} else {
			log.Printf("error getting project by id: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	err = h.writeJSON(w, http.StatusOK, project, nil)
	if err != nil {
		log.Printf("error writing get project response: %v", err)
	}
}

func (h *Handlers) GetAllProjectsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	projects, err := h.projectStore.GetAll(userID)
	if err != nil {
		log.Printf("error getting all projects: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.writeJSON(w, http.StatusOK, projects, nil)
	if err != nil {
		log.Printf("error writing get all projects response: %v", err)
	}
}

func (h *Handlers) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "projectID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Bad Request: invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := h.projectStore.GetByID(id, userID)
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
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	err = h.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Bad Request: could not decode JSON", http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		project.Name = *input.Name
	}

	if input.Description != nil {
		project.Description = input.Description
	}

	err = h.projectStore.Update(project)
	if err != nil {
		log.Printf("error updating project: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.writeJSON(w, http.StatusOK, project, nil)
	if err != nil {
		log.Printf("error writing update project response: %v", err)
	}
}

func (h *Handlers) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "projectID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Bad Request: invalid project ID", http.StatusBadRequest)
		return
	}

	err = h.projectStore.Delete(id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not Found: project not found", http.StatusNotFound)
		} else {
			log.Printf("error deleting project: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
