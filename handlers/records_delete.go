package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

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
	payload := IDPayload{ID: rid}
	json.NewEncoder(w).Encode(payload)
}
