package handlers

import (
	"balancer-api/balancer"
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

type Handler struct {
	RecordService balancer.RecordService
}

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
	payload := idPayload{ID: *rid}
	json.NewEncoder(w).Encode(payload)
}

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

	if newRecord.Name != "" {
		record.Name = newRecord.Name
	}

	if newRecord.Type == Asset || newRecord.Type == Liability {
		record.Type = strings.ToUpper(newRecord.Type)
	}

	record.Balance = newRecord.Balance

	err = h.RecordService.UpdateRecord(record)
	if err != nil {
		http.Error(w, "Error reading entry from db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	recordID := chi.URLParam(r, "id")

	rid, err := strconv.ParseUint(recordID, 10, 64)
	if err != nil {
		log.Println("Error parsing uint ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.RecordService.DeleteRecord(rid); err != nil {
		http.Error(w, "Error reading entries in db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	payload := idPayload{ID: rid}
	json.NewEncoder(w).Encode(payload)
}

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
	payload := valuePayload{Value: sum}
	json.NewEncoder(w).Encode(payload)
}

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
	payload := valuePayload{Value: sum}
	json.NewEncoder(w).Encode(payload)
}
