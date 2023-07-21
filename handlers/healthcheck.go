package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/olad5/go-url-shortener/config"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	type healthCheckResponse struct {
		Status string `json:"status"`
	}

	healthStatus := config.RepositoryAdapter.Ping()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(healthCheckResponse{Status: string(healthStatus)}); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}
