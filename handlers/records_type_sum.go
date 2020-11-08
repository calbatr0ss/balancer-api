package handlers

import (
	"balancer-api/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) GetTypeSum(w http.ResponseWriter, r *http.Request) {
	var records *[]models.Record
	var sum float64

	keys, ok := r.URL.Query()["type"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'type' is missing")
		http.Error(w, "type not included in request", http.StatusBadRequest)
		return
	}

	typeKey := strings.ToUpper(keys[0])

	// Ensure provided type matches a valid type
	if typeKey != Asset && typeKey != Liability {
		log.Println("Url Param 'type' is invalid")
		http.Error(w, "invalid type provided", http.StatusBadRequest)
		return
	}

	records, err := h.RecordService.GetRecordsByType(typeKey)
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
