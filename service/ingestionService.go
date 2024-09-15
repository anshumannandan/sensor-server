package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"net/http"
	"sensor-server/dataAccess"
)

func IngestionService(csvURL string) error {
	log.Printf("Ingesting sensor data from: %s", csvURL)

	resp, err := http.Get(csvURL)
	if err != nil {
		log.Printf("Failed to fetch CSV from %s: %v", csvURL, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("received non-OK HTTP status")
		log.Printf("Failed to fetch CSV: %s, status code: %d", resp.Status, resp.StatusCode)
		return err
	}

	csvData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read CSV data: %v", err)
		return err
	}

	records, err := csv.NewReader(bytes.NewReader(csvData)).ReadAll()
	if err != nil {
		log.Printf("Failed to parse CSV data: %v", err)
		return err
	}

	if len(records) == 0 {
		err := errors.New("no records found in CSV data")
		log.Println(err)
		return err
	}

	err = dataAccess.IngestSensorData(records)
	if err != nil {
		log.Printf("Failed to ingest sensor data: %v", err)
		return err
	}

	log.Printf("Successfully ingested %d records", len(records))
	return nil
}
