package main

import (
	"log"
	"net/http"
	"sensor-server/controller"
	"sensor-server/initializer"
)

func main() {
	dbURL := "http://influxdb:8086"
	dbToken := "admin-token"
	dbOrg := "sensor-org"
	dbBucket := "sensor-data-bucket"

    if err := initializer.Initialize(dbURL, dbToken, dbOrg, dbBucket); err != nil {
        log.Fatalf("Initialization failed: %v", err)
    } else {
		log.Println("Initialization successful")
	}

	http.HandleFunc("/ingest", controller.IngestionController);
	http.HandleFunc("/median", controller.RetrievalController);

	log.Println("Server Listening on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
