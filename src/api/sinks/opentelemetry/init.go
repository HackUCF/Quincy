package opentelemetry

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otellog "go.opentelemetry.io/otel/log"
	otelsdklog "go.opentelemetry.io/otel/sdk/log"
)

var logger otellog.Logger

// Init sets up the OTLP HTTP log exporter.
// Scheme determines TLS: http → insecure, https → TLS.
// Set Username or BasicAuth to enable HTTP Basic auth.
func Init(ctx context.Context, cfg config.OTelConfig) (func(), error) {
	opts, err := buildOpts(cfg)
	if err != nil {
		return nil, err
	}

	exp, err := otlploghttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	// default batch size options
	batchSize := cfg.Batching.BatchSize
	if batchSize == 0 {
		batchSize = 20
	}
	exportIntervalSecs := cfg.Batching.ExportInterval
	if exportIntervalSecs == 0 {
		exportIntervalSecs = 5
	}
	exportInterval := time.Duration(exportIntervalSecs) * time.Second
	maxQueueSize := cfg.Batching.MaxQueueSize
	if maxQueueSize == 0 {
		maxQueueSize = 200
	}

	// configure batching
	processor := otelsdklog.NewBatchProcessor(exp,
		otelsdklog.WithExportMaxBatchSize(batchSize),
		otelsdklog.WithExportInterval(exportInterval),
		otelsdklog.WithMaxQueueSize(maxQueueSize),
	)

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Error("otel sink error", "error", err)
	}))

	provider := otelsdklog.NewLoggerProvider(otelsdklog.WithProcessor(processor))
	logger = provider.Logger("quincy")

	return func() { provider.Shutdown(context.Background()) }, nil
}

// buildOpts constructs the OTLP HTTP exporter options from cfg.
func buildOpts(cfg config.OTelConfig) ([]otlploghttp.Option, error) {
	u, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid OTel endpoint %q: %w", cfg.Endpoint, err)
	}

	opts := []otlploghttp.Option{
		otlploghttp.WithEndpoint(u.Host),
	}

	// allow http endpoints
	if u.Scheme == "http" {
		opts = append(opts, otlploghttp.WithInsecure())
	}

	// use url paths
	if u.Path != "" && u.Path != "/" {
		opts = append(opts, otlploghttp.WithURLPath(u.Path))
	}

	headers := map[string]string{}

	// use or create basic auth headers
	if cfg.BasicAuth != "" {
		headers["Authorization"] = "Basic " + cfg.BasicAuth
	} else if cfg.Username != "" {
		creds := base64.StdEncoding.EncodeToString([]byte(cfg.Username + ":" + cfg.Password))
		headers["Authorization"] = "Basic " + creds
	}

	// set open observe stream
	if cfg.StreamName != "" {
		headers["stream-name"] = cfg.StreamName
	}

	// add headers if any are set
	if len(headers) > 0 {
		opts = append(opts, otlploghttp.WithHeaders(headers))
	}

	return opts, nil
}
