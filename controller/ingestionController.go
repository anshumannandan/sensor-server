package controller

import (
	"log"
	"net/http"
	"net/url"
	"sensor-server/service"
)

func IngestionController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	csvURL := r.URL.Query().Get("url")

	decodedURL, err := url.QueryUnescape(csvURL)
	if err != nil {
		log.Printf("Error decoding URL: %v", err)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(decodedURL)
	if err != nil {
		log.Printf("Invalid URL format: %v", err)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	err = service.IngestionService(decodedURL)
	if err != nil {
		log.Printf("Error processing ingestion: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Ingestion successful"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
