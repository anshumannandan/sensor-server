package controller

import (
	"net/http"
	"sensor/service"
)

func IngestionController(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

	csvURL := r.URL.Query().Get("url")
    if csvURL == "" {
        http.Error(w, "Missing CSV URL", http.StatusBadRequest)
        return
    }

	try {
		service.IngestionService(csvURL)
	} catch (e) {
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Ingestion successful"))
}
