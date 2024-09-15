package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"sensor-server/service"
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
			log.Printf("Error unmarshaling JSON: %v", err)
			http.Error(w, "Invalid filter format", http.StatusBadRequest)
			return
		}
	}

	median, err := service.RetrievalService(decodedData)
	if err != nil {
		log.Printf("Error retrieving median: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		// "count" : count,
		"median": median,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
