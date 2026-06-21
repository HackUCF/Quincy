package opentelemetry

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/types"
	testcontainer "github.com/HackUCF/quincy/testutil/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// collectorCfg is a minimal OTel Collector config: OTLP HTTP receiver + debug exporter.
// The debug exporter with detailed verbosity logs every received record to stdout,
// which lets us assert delivery by reading container logs.
const collectorCfg = `
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
exporters:
  debug:
    verbosity: detailed
service:
  pipelines:
    logs:
      receivers: [otlp]
      exporters: [debug]
`

func startCollector(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	t.Helper()

	cfgPath := filepath.Join(t.TempDir(), "otel-config.yaml")
	if err := os.WriteFile(cfgPath, []byte(collectorCfg), 0o644); err != nil {
		t.Logf("WARNING: failed to write collector config, skipping %s: %v", t.Name(), err)
		t.Skipf("collector config write failed: %v", err)
	}

	c := testcontainer.RunContainer(ctx, t, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "otel/opentelemetry-collector-contrib:latest",
			ExposedPorts: []string{"4318/tcp"},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      cfgPath,
					ContainerFilePath: "/etc/otelcol-contrib/config.yaml",
					FileMode:          0o644,
				},
			},
			WaitingFor: wait.ForLog("Everything is ready").
				WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})

	host, err := c.Host(ctx)
	if err != nil {
		t.Logf("WARNING: could not get container host, skipping %s: %v", t.Name(), err)
		t.Skipf("container host unavailable: %v", err)
	}
	port, err := c.MappedPort(ctx, "4318/tcp")
	if err != nil {
		t.Logf("WARNING: could not get mapped port, skipping %s: %v", t.Name(), err)
		t.Skipf("container port unavailable: %v", err)
	}

	return c, "http://" + host + ":" + port.Port()
}

func TestIntegration_scoreReachesCollector(t *testing.T) {
	ctx := context.Background()
	c, endpoint := startCollector(ctx, t)

	cfg := config.OTelConfig{
		Endpoint: endpoint,
		Batching: config.OTelBatching{
			BatchSize:      1,
			ExportInterval: 1,
			MaxQueueSize:   10,
		},
	}

	shutdown, err := Init(ctx, cfg)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer shutdown()

	score := types.Score{
		BoxName:     "integration-box",
		ServiceName: "integration-svc",
		TeamNum:     42,
		Status:      true,
		Message:     "integration test ok",
	}
	if err := AddScore(ctx, score); err != nil {
		t.Fatalf("AddScore: %v", err)
	}

	// Wait for the batch interval to flush + collector to write logs.
	time.Sleep(3 * time.Second)

	logs, err := c.Logs(ctx)
	if err != nil {
		t.Fatalf("container logs: %v", err)
	}
	defer logs.Close()

	var buf strings.Builder
	io.Copy(&buf, logs)
	output := buf.String()

	for _, want := range []string{"integration-box", "integration-svc", "integration test ok"} {
		if !strings.Contains(output, want) {
			t.Errorf("collector logs missing %q\nfull output:\n%s", want, output)
		}
	}
}

func TestIntegration_failedCheckReachesCollector(t *testing.T) {
	ctx := context.Background()
	c, endpoint := startCollector(ctx, t)

	cfg := config.OTelConfig{
		Endpoint: endpoint,
		Batching: config.OTelBatching{
			BatchSize:      1,
			ExportInterval: 1,
			MaxQueueSize:   10,
		},
	}

	shutdown, err := Init(ctx, cfg)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer shutdown()

	if err := AddScore(ctx, types.Score{
		BoxName:     "fail-box",
		ServiceName: "fail-svc",
		TeamNum:     7,
		Status:      false,
		Message:     "connection refused",
	}); err != nil {
		t.Fatalf("AddScore: %v", err)
	}

	time.Sleep(3 * time.Second)

	logs, err := c.Logs(ctx)
	if err != nil {
		t.Fatalf("container logs: %v", err)
	}
	defer logs.Close()

	var buf strings.Builder
	io.Copy(&buf, logs)
	output := buf.String()

	for _, want := range []string{"fail-box", "fail-svc", "connection refused"} {
		if !strings.Contains(output, want) {
			t.Errorf("collector logs missing %q\nfull output:\n%s", want, output)
		}
	}
}
