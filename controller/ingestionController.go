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
	if csvURL == "" {
		log.Printf("url parameter is required")
		http.Error(w, "url parameter is required", http.StatusBadRequest)
		return
	}

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

	if err = service.IngestionService(decodedURL); err != nil {
		log.Printf("Error processing ingestion: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write([]byte("Ingestion successful")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
