package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"sensor-server/service"
	"strconv"
)

func RetrievalController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Method Not Allowed: %v", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	filterJSON := r.URL.Query().Get("filter")
	var decodedData map[string][]string

	if filterJSON != "" {
		err := json.Unmarshal([]byte(filterJSON), &decodedData)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
		}
	}

	median, err := service.RetrievalService(decodedData)
	if err != nil {
		log.Printf("Error retrieving median: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	medianStr := strconv.Itoa(int(median))
	_, err = w.Write([]byte(medianStr))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
