package controller

import (
	"net/http"
	"sensor/service"
)

func RetrievalController(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	try {
		median, err := service.RetrievalService()
	} catch (e) {
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(median))
}
