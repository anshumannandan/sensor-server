package service

import (
	"context"
	"fmt"
	"log"
	"sensor-server/initializer"
	"strings"
)

func RetrievalService(filter map[string][]string) (float64, error) {
	log.Printf("Retrieving sensor data with filter: %v", filter)

	var builder strings.Builder
	builder.WriteString(`from(bucket: "sensor-data-bucket") |> range(start: 0)`)

	var conditions []string
	if ids, ok := filter["id"]; ok && len(ids) > 0 {
		conditions = append(conditions, buildFilterCondition("id", ids))
	}
	if types, ok := filter["type"]; ok && len(types) > 0 {
		conditions = append(conditions, buildFilterCondition("type", types))
	}
	if subtypes, ok := filter["subtype"]; ok && len(subtypes) > 0 {
		conditions = append(conditions, buildFilterCondition("subtype", subtypes))
	}
	if locations, ok := filter["location"]; ok && len(locations) > 0 {
		conditions = append(conditions, buildFilterCondition("location", locations))
	}

	if len(conditions) > 0 {
		builder.WriteString(` |> filter(fn: (r) => `)
		builder.WriteString(strings.Join(conditions, " and "))
		builder.WriteString(")")
	}

	builder.WriteString(` |> group(columns: ["_measurement"]) |> median()`)

	query := builder.String()
	log.Printf("InfluxDB Query: %s", query)

	result, err := initializer.QueryAPI.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying InfluxDB: %v", err)
		return 0, err
	}
	defer result.Close()

	var median float64
	if result.Next() {
		record := result.Record()
		log.Printf("Record: %v", record)
		if value, ok := record.Value().(float64); ok {
			median = value
		} else {
			log.Printf("Error parsing median value")
		}
	} else {
		if result.Err() != nil {
			log.Printf("Error retrieving data: %v", result.Err())
			return median, result.Err()
		}
		log.Printf("No records found")
	}
	
	log.Printf("Retrieved median value: %f", median)
	return median, nil
}

func buildFilterCondition(field string, values []string) string {
	conditions := make([]string, len(values))
	for i, value := range values {
		conditions[i] = fmt.Sprintf(`r.%s == "%s"`, field, value)
	}
	return fmt.Sprintf(`(%s)`, strings.Join(conditions, " or "))
}
