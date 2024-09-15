package dataAccess

import (
	"context"
	"fmt"
	"log"
	"sensor-server/initializer"
	"strings"
)

func RetrieveSensorData(filter map[string][]string) (float64, error) {
	query := `from(bucket: "sensor-data-bucket") |> range(start: 0)`

	var conditions []string
	if ids, ok := filter["id"]; ok && len(ids) > 0 {
		conditions = append(conditions, fmt.Sprintf(`(r.id == "%s")`, strings.Join(ids, `" or r.id == "`)))
	}
	if types, ok := filter["type"]; ok && len(types) > 0 {
		conditions = append(conditions, fmt.Sprintf(`(r.type == "%s")`, strings.Join(types, `" or r.type == "`)))
	}
	if subtypes, ok := filter["subtype"]; ok && len(subtypes) > 0 {
		conditions = append(conditions, fmt.Sprintf(`(r.subtype == "%s")`, strings.Join(subtypes, `" or r.subtype == "`)))
	}
	if locations, ok := filter["location"]; ok && len(locations) > 0 {
		conditions = append(conditions, fmt.Sprintf(`(r.location == "%s")`, strings.Join(locations, `" or r.location == "`)))
	}

	if len(conditions) > 0 {
		query += fmt.Sprintf(` |> filter(fn: (r) => %s)`, strings.Join(conditions, " and "))
	}

	query += ` |> group(columns : ["reading"]) |> median()`
	log.Printf("InfluxDB Query: %s", query)

	result, err := initializer.QueryAPI.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying InfluxDB: %v", err)
		return 0, err
	}
	defer result.Close()

	var median float64
	for result.Next() {
		record := result.Record()
		log.Printf("Record: %v", record)
		if value, ok := record.Value().(float64); ok {
			median = value
		} else {
			log.Printf("Error parsing median value")
		}
	}
	if result.Err() != nil {
		log.Printf("Error retrieving data: %v", result.Err())
		return median, result.Err()
	}

	return median, nil
}
