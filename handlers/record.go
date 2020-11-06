package handlers

import (
	"balancer-api/db"
	"balancer-api/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type idPayload struct {
	ID uint64 `json:"recordId"`
}

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	var records []models.Record

	if result := db.DB.Find(&records); result.Error != nil {
		http.Error(w, "Error reading entries in db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(records)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
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

	if result := db.DB.Create(&record); result.Error != nil {
		http.Error(w, "Error creating entry in db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	payload := idPayload{ID: record.ID}
	json.NewEncoder(w).Encode(payload)
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	var record models.Record
	var newRecord models.Record
	recordID := chi.URLParam(r, "id")

	if result := db.DB.Find(&record, recordID); result.Error != nil {
		http.Error(w, "Error reading entry in db", http.StatusInternalServerError)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&newRecord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newRecord.Name != "" {
		record.Name = newRecord.Name
	}

	if newRecord.Type != "" {
		record.Type = newRecord.Type
	}

	record.Balance = newRecord.Balance

	if result := db.DB.Save(&record); result.Error != nil {
		http.Error(w, "Error updating entry in db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	recordID := chi.URLParam(r, "id")

	if result := db.DB.Delete(&models.Record{}, recordID); result.Error != nil {
		http.Error(w, "Error deleting entry in db", http.StatusInternalServerError)
		return
	}

	rid, err := strconv.ParseUint(recordID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	payload := idPayload{ID: rid}
	json.NewEncoder(w).Encode(payload)
}
