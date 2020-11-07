package handlers

import (
	"balancer-api/db"
	"balancer-api/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

const (
	Asset     = "ASSET"
	Liability = "LIABILITY"
)

type idPayload struct {
	ID uint64 `json:"recordId"`
}

type valuePayload struct {
	Value float64 `json:"value"`
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

	// Ensure the type is valid
	cleanType := strings.ToUpper(record.Type)
	if cleanType != Asset && cleanType != Liability {
		log.Println("record type is invalid")
		http.Error(w, "invalid type provided", http.StatusBadRequest)
		return
	}

	// Set the cleaned type
	record.Type = cleanType

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
		log.Println("Error reading entries from db")
		http.Error(w, "Error reading entries from db", http.StatusInternalServerError)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&newRecord)
	if err != nil {
		log.Println("Error decoding json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newRecord.Name != "" {
		record.Name = newRecord.Name
	}

	if newRecord.Type == Asset || newRecord.Type == Liability {
		record.Type = strings.ToUpper(newRecord.Type)
	}

	record.Balance = newRecord.Balance

	if result := db.DB.Save(&record); result.Error != nil {
		log.Println("Error updating entry in db")
		http.Error(w, "Error updating entry in db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	recordID := chi.URLParam(r, "id")

	if result := db.DB.Delete(&models.Record{}, recordID); result.Error != nil {
		log.Println("Error deleting entry in db")
		http.Error(w, "Error deleting entry in db", http.StatusInternalServerError)
		return
	}

	rid, err := strconv.ParseUint(recordID, 10, 64)
	if err != nil {
		log.Println("Error parsing uint ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	payload := idPayload{ID: rid}
	json.NewEncoder(w).Encode(payload)
}

func GetNetWorth(w http.ResponseWriter, r *http.Request) {
	var records []models.Record
	var sum float64

	if result := db.DB.Find(&records); result.Error != nil {
		log.Println("Error reading entries in db")
		http.Error(w, "Error reading entries in db", http.StatusInternalServerError)
		return
	}

	for _, r := range records {
		sum += r.Balance
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	payload := valuePayload{Value: sum}
	json.NewEncoder(w).Encode(payload)
}

func GetTypeSum(w http.ResponseWriter, r *http.Request) {
	var records []models.Record
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

	if result := db.DB.Where("type = ?", typeKey).Find(&records); result.Error != nil {
		log.Println("Error reading entries in db")
		http.Error(w, "Error reading entries in db", http.StatusInternalServerError)
		return
	}

	for _, r := range records {
		sum += r.Balance
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	payload := valuePayload{Value: sum}
	json.NewEncoder(w).Encode(payload)
}
