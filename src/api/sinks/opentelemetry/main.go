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
	// the check's stdout is the record body; stderr rides along as its own
	// attribute so the two streams stay separable downstream.
	r.SetBody(attribute.StringValue(score.Stdout))
	r.SetSeverity(log.SeverityInfo)
	r.AddAttributes(
		attribute.Int64("team", int64(score.TeamNum)),
		attribute.String("box", string(score.BoxName)),
		attribute.String("service", string(score.ServiceName)),
		attribute.Bool("status", score.Status),
		attribute.String("stderr", score.Stderr),
	)

	logger.Emit(ctx, r)
	return nil
}
