package controller

import (
	"net/http"
	"sensor-server/service"
	"strconv"
)

func RetrievalController(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	median, err := service.RetrievalService()
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.WriteHeader(http.StatusOK)
	medianStr := strconv.FormatFloat(median, 'f', -1, 64)
	w.Write([]byte(medianStr))
}
