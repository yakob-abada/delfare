package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/yakob-abada/delfare/reader-service/config"
	"github.com/yakob-abada/delfare/reader-service/domain"
)

type InfluxDBRepository struct {
	client influxdb2.Client
	org    string
	bucket string
	logger domain.Logger
}

func NewInfluxDBRepository(cfg config.Config, logger domain.Logger) domain.EventRepository {
	client := influxdb2.NewClient(cfg.InfluxDBURL, cfg.InfluxDBToken)
	return &InfluxDBRepository{client, cfg.InfluxDBOrg, cfg.InfluxDBBucket, logger}
}

func (r *InfluxDBRepository) GetCriticalEvents(ctx context.Context, limit, criticality int) ([]domain.Event, error) {
	query := fmt.Sprintf(`from(bucket: "%s")
		  |> range(start: -30d) 
		  |> filter(fn: (r) => r._measurement == "security_events")
		  |> filter(fn: (r) => r._field == "criticality" or r._field == "message")
		  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
		  |> filter(fn: (r) => r.criticality > 5) // Filter events with criticality > %d
		  |> sort(columns: ["_time"], desc: true) // Sort by time (latest first)
		  |> limit(n: %d) // Get the latest 100 events`, r.bucket, criticality, limit)

	queryAPI := r.client.QueryAPI(r.org)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		r.logger.Error(domain.LogContext{}, "InfluxDB query failed", err)
		return nil, err
	}

	var events []domain.Event
	for result.Next() {
		events = append(events, domain.Event{
			RequestID:   fmt.Sprintf("%s", result.Record().ValueByKey("request_id")),
			Criticality: int(result.Record().ValueByKey("criticality").(int64)),
			Timestamp:   result.Record().Time(),
			Message:     result.Record().ValueByKey("message").(string),
		})
	}

	return events, nil
}
