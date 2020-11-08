package handlers

import (
	"balancer-api/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	var newRecord models.Record
	recordID := chi.URLParam(r, "id")

	record, err := h.RecordService.GetRecord(recordID)
	if err != nil {
		http.Error(w, "Error reading entry from db", http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&newRecord)
	if err != nil {
		log.Println("Error decoding json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure name and type are included
	if newRecord.Name == "" || newRecord.Type == "" {
		http.Error(w, "name or type not included in request body", http.StatusBadRequest)
		return
	}

	// Ensure the type is valid
	cleanType := strings.ToUpper(newRecord.Type)
	if cleanType != Asset && cleanType != Liability {
		log.Println("record type is invalid")
		http.Error(w, "invalid type provided", http.StatusBadRequest)
		return
	}

	// Set the cleaned type
	newRecord.Type = cleanType

	// Set new record with ID
	newRecord.ID = record.ID

	err = h.RecordService.UpdateRecord(&newRecord)
	if err != nil {
		http.Error(w, "Error reading entry from db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
