package handlers

import (
	"balancer-api/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	var newRecord models.Record
	recordID := chi.URLParam(r, "id")
	fmt.Println("here")

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
	fmt.Println("here")

	if newRecord.Name != "" {
		record.Name = newRecord.Name
	}

	if newRecord.Type == Asset || newRecord.Type == Liability {
		record.Type = strings.ToUpper(newRecord.Type)
	}

	err = h.RecordService.UpdateRecord(record)
	if err != nil {
		http.Error(w, "Error reading entry from db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
