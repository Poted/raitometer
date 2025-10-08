package handlers

import (
	"log"
	"net/http"
)

type healthCheckResponse struct {
	APIVersion string `json:"apiVersion"`
	APISystem  string `json:"apiSystem"`
	Database   string `json:"database"`
}

func (h *Handlers) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	dbStatus := "OK"
	err := h.db.Ping()
	if err != nil {
		dbStatus = "Unavailable"
		log.Printf("warning: healthcheck detected database issue: %v", err)
	}

	response := healthCheckResponse{
		APIVersion: "1.0.0",
		APISystem:  "raitometer Core API",
		Database:   dbStatus,
	}

	err = h.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		log.Printf("error writing healthcheck response: %v", err)
	}
}
