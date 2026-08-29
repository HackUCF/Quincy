package opentelemetry

import (
	"context"
	"time"

	"github.com/HackUCF/quincy/common/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

// AddScore emits a log record for a completed score check.
func AddScore(ctx context.Context, score types.Score) error {

	var r log.Record
	r.SetTimestamp(time.UnixMicro(score.Timestamp))
	r.SetBody(attribute.StringValue(score.Message))
	r.SetSeverity(log.SeverityInfo)
	r.AddAttributes(
		attribute.Int64("team", int64(score.TeamNum)),
		attribute.String("box", string(score.BoxName)),
		attribute.String("service", string(score.ServiceName)),
		attribute.Bool("status", score.Status),
	)

	logger.Emit(ctx, r)
	return nil
}
