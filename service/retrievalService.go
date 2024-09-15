package service

import (
	"context"
	"fmt"
	"log"
	"sensor-server/initializer"
	"strings"
)

func RetrievalService(filter map[string][]string) (float64, int64, error) {
	log.Printf("Retrieving sensor data with filter: %v", filter)

	query := buildQuery(filter)
	log.Printf("InfluxDB Query: %s", query)

	result, err := initializer.QueryAPI.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying InfluxDB: %v", err)
		return 0, 0, err
	}
	defer result.Close()

	var (
		median float64
		count  int64
	)

	for result.Next() {
		record := result.Record()
		log.Printf("Record: %v", record)
		value := record.Value()
		resultName := record.ValueByKey("result")

		switch resultName {
		case "median":
			if medianValue, ok := value.(float64); ok {
				median = medianValue
			} else {
				log.Printf("Error parsing median value")
			}
		case "count":
			if countValue, ok := value.(int64); ok {
				count = countValue
			} else {
				log.Printf("Error parsing count value")
			}
		}
	}

	if err := result.Err(); err != nil {
		log.Printf("Error retrieving data: %v", err)
	}

	log.Printf("Retrieved median: %f, count: %d", median, count)
	return median, count, nil
}

func buildQuery(filter map[string][]string) string {
	var builder strings.Builder
	builder.WriteString(`data = from(bucket: "sensor-data-bucket") |> range(start: 0)`)

	var conditions []string
	if types, ok := filter["type"]; ok && len(types) > 0 {
		conditions = append(conditions, buildFilterCondition("type", types))
	}
	if subtypes, ok := filter["subtype"]; ok && len(subtypes) > 0 {
		conditions = append(conditions, buildFilterCondition("subtype", subtypes))
	}
	appendFilterCondition(&builder, conditions)
	conditions = []string{}

	if locations, ok := filter["location"]; ok && len(locations) > 0 {
		conditions = append(conditions, buildFilterCondition("location", locations))
	}
	appendFilterCondition(&builder, conditions)
	conditions = []string{}

	if ids, ok := filter["id"]; ok && len(ids) > 0 {
		conditions = append(conditions, buildFilterCondition("id", ids))
	}
	appendFilterCondition(&builder, conditions)

	builder.WriteString(` |> group(columns: ["_measurement"]) `)
	builder.WriteString(` data |> median() |> yield(name: "median") `)
	builder.WriteString(` data |> count() |> yield(name: "count") `)
	return builder.String()
}

func buildFilterCondition(field string, values []string) string {
	conditions := make([]string, len(values))
	for i, value := range values {
		conditions[i] = fmt.Sprintf(`r.%s == "%s"`, field, value)
	}
	return fmt.Sprintf(`(%s)`, strings.Join(conditions, " or "))
}

func appendFilterCondition(builder *strings.Builder, conditions []string) {
	if len(conditions) > 0 {
		builder.WriteString(` |> filter(fn: (r) => `)
		builder.WriteString(strings.Join(conditions, " and "))
		builder.WriteString(")")
	}
}
