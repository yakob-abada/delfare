package infrastructure

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/yakob-abada/delfare/writer-service/config"
	"github.com/yakob-abada/delfare/writer-service/domain"
)

type InfluxDBRepository struct {
	client influxdb2.Client
	org    string
	bucket string
	logger domain.Logger
}

func NewInfluxDBRepository(cfg config.Config, logger domain.Logger) domain.Repository {
	client := influxdb2.NewClient(cfg.InfluxDBURL, cfg.InfluxDBToken)
	return &InfluxDBRepository{client, cfg.InfluxDBOrg, cfg.InfluxDBBucket, logger}
}

// Write writes a security event to InfluxDB
func (r *InfluxDBRepository) Write(ctx context.Context, event domain.Event) error {
	writeAPI := r.client.WriteAPIBlocking(r.org, r.bucket)

	p := influxdb2.NewPointWithMeasurement("security_events").
		AddTag("request_id", event.RequestID).
		AddField("criticality", event.Criticality).
		AddField("message", event.Message).
		SetTime(event.Timestamp)

	err := writeAPI.WritePoint(ctx, p)
	if err != nil {
		r.logger.Error(domain.LogContext{}, "Error writing to InfluxDB", "error", err)
		return err
	}
	return nil
}
