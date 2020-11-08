package handlers

import (
	"balancer-api/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	var record models.Record

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if record.Name == "" || record.Type == "" {
		http.Error(w, "name or type not included in request body", http.StatusBadRequest)
		return
	}

	// Ensure the type is valid
	cleanType := strings.ToUpper(record.Type)
	if cleanType != Asset && cleanType != Liability {
		log.Println("record type is invalid")
		http.Error(w, "invalid type provided", http.StatusBadRequest)
		return
	}

	// Set the cleaned type
	record.Type = cleanType

	rid, err := h.RecordService.CreateRecord(&record)
	if err != nil {
		http.Error(w, "Error reading entries in db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	payload := IDPayload{ID: *rid}
	json.NewEncoder(w).Encode(payload)
}
