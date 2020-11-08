package handlers

import (
	"balancer-api/models"
	"encoding/json"
	"net/http"
)

func (h *Handler) GetNetWorth(w http.ResponseWriter, r *http.Request) {
	var records *[]models.Record
	var sum float64

	records, err := h.RecordService.GetAllRecords()
	if err != nil {
		http.Error(w, "Error reading entries in db", http.StatusInternalServerError)
		return
	}

	for _, r := range *records {
		sum += r.Balance
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	payload := ValuePayload{Value: sum}
	json.NewEncoder(w).Encode(payload)
}
