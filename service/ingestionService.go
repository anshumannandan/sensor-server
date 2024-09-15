package service

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"net/http"
	"sensor-server/initializer"
	"strconv"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
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
		err := errors.New("received non-OK HTTP status: " + resp.Status)
		log.Printf("Failed to fetch CSV: %s, status code: %d", resp.Status, resp.StatusCode)
		return err
	}

	reader := csv.NewReader(resp.Body)
	if err := IngestSensorData(reader); err != nil {
		log.Printf("Failed to ingest sensor data: %v", err)
		return err
	}
	return nil
}

func IngestSensorData(reader *csv.Reader) error {
	header := true
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading CSV record: %v", err)
			continue
		}
		if header {
			header = false
			continue
		}

		point, err := createPoint(record)
		if err != nil {
			log.Printf("Error creating point: %v", err)
			continue
		}
		initializer.WriteAPI.WritePoint(point)
	}
	initializer.WriteAPI.Flush()
	return nil
}

func createPoint(record []string) (*write.Point, error) {
	timestamp, err := time.Parse("2006-01-02 15:04:05", record[5])
	if err != nil {
		return nil, err
	}

	reading, err := strconv.Atoi(record[3])
	if err != nil {
		return nil, err
	}

	point := influxdb2.NewPointWithMeasurement("sensor-data-measurement").
		AddTag("id", record[0]).
		AddTag("type", record[1]).
		AddTag("subtype", record[2]).
		AddField("reading", reading).
		AddTag("location", record[4]).
		SetTime(timestamp)

	return point, nil
}
