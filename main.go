package main

import (
	"net/http"
	"sensor/controller"
)

func main() {
	http.HandleFunc("/ingest", controller.IngestionController);
	http.HandleFunc("/median", controller.RetrievalController);

	http.ListenAndServe(":5000", nil)
}
