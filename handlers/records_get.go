package handlers

import (
	"balancer-api/models"
	"encoding/json"
	"net/http"
)

func (h *Handler) GetAllRecords(w http.ResponseWriter, r *http.Request) {
	var records *[]models.Record

	records, err := h.RecordService.GetAllRecords()
	if err != nil {
		http.Error(w, "Error reading entries in db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(records)
}
