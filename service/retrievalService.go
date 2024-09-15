package service

import (
	"log"
	"sensor-server/dataAccess"
)

func RetrievalService(data map[string][]string) (int, error) {
	log.Printf("Retrieving sensor data with filter: %v", data)
	median, err := dataAccess.RetrieveSensorData(data)
	if err != nil {
		log.Printf("Error retrieving sensor data: %v", err)
		return 0, err
	}
	
	medianInt := int(median)
	log.Printf("Retrieved median value: %d", medianInt)
	return medianInt, nil
}
