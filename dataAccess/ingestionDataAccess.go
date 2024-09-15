package dataAccess

import (
	// "context"
	// "log"
	// "sensor-server/initializer"
	// "time"

	// "github.com/influxdata/influxdb-client-go/v2"
	// "github.com/influxdata/influxdb-client-go/v2/api/write"
)

func IngestSensorData(records [][]string) error {
	// points := write.NewBatchPoint()

	// for _,row := range records[1:] {
	// 	timestamp, err := time.Parse(time.RFC3339, row[5])
	// 	if err != nil {
	// 		log.Printf("Error parsing timestamp: %v", err)
	// 	}
	// 	p := influxdb2.NewPointWithMeasurement("sensor-data-measurement").
	// 		AddTag("id", row[0]).
	// 		AddTag("type", row[1]).
	// 		AddTag("subtype", row[2]).
	// 		AddField("reading", row[3]).
	// 		AddTag("location", row[4]).
	// 		SetTime(timestamp)

	// 	points = append(points, p)
	// }

	// err := initializer.WriteAPI.Write(context.Background(), points)
    // if err != nil {
	// 	log.Printf("Error writing points: %v", err)
	// 	return err
    // }

	return nil
}
