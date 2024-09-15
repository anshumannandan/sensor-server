package service

import (
	"log"
	"sensor-server/dataAccess"
)

func RetrievalService(data map[string][]string) (int, error) {
	log.Printf("Retrieving sensor data with filter: %v", data)
	median, err := dataAccess.RetrieveSensorData(data)
	medianInt := int(median)
	return medianInt, err
}
