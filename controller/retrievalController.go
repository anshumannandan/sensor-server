package controller

import (
	"net/http"
	"log"
	"strconv"
	"sensor-server/service"
)

func RetrievalController(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		log.Printf("Method Not Allowed: %v", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	median, err := service.RetrievalService()
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
