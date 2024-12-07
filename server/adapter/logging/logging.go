package logging

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"
)

func SendLogToInfluxDB(message string, level string, influxDBClient influxdb2.Client, org string, bucket string) error {
	writeAPI := influxDBClient.WriteAPIBlocking(org, bucket)
	tags := map[string]string{"level": level}
	fields := map[string]interface{}{
		"message": message,
	}
	p := influxdb2.NewPoint("application_logs",
		tags,
		fields,
		time.Now())
	return writeAPI.WritePoint(context.Background(), p)
}
