package handlers

import (
	"bytes"
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) AnalyzeProjectImageHandler(w http.ResponseWriter, r *http.Request) {
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

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Bad Request: could not get image from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		log.Printf("error creating form file for AI service: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Printf("error copying file data for AI service: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	writer.Close()

	aiServiceURL := "http://raitometer_ai:8000/analyze-image"
	req, err := http.NewRequest("POST", aiServiceURL, &requestBody)
	if err != nil {
		log.Printf("error creating request to AI service: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending request to AI service: %v", err)
		http.Error(w, "Service Unavailable: AI service is not reachable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
