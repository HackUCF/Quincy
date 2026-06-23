package opentelemetry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/types"
)

// fastBatchCfg returns a base config pointing at endpoint with aggressive flush settings for tests.
func fastBatchCfg(endpoint string) config.OTelConfig {
	return config.OTelConfig{
		Endpoint: endpoint,
		Batching: config.OTelBatching{
			BatchSize:      1,
			ExportInterval: 1, // 1 second — minimum granularity, fast enough for tests
			MaxQueueSize:   10,
		},
	}
}

// captureServer starts a test server that captures the first incoming request's headers.
// The returned channel receives the captured request; done is closed after the first call.
func captureServer() (*httptest.Server, <-chan *http.Request) {
	ch := make(chan *http.Request, 1)
	var once sync.Once
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() { ch <- r.Clone(context.Background()) })
		w.WriteHeader(http.StatusOK)
	}))
	return srv, ch
}

func TestBuildOpts_invalidURL(t *testing.T) {
	_, err := buildOpts(config.OTelConfig{Endpoint: "://bad"})
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestBuildOpts_basicAuth_preencoded(t *testing.T) {
	srv, ch := captureServer()
	defer srv.Close()

	cfg := fastBatchCfg(srv.URL)
	cfg.BasicAuth = "dXNlcjpwYXNz" // base64("user:pass")
	cfg.Username = "ignored"
	cfg.Password = "ignored"

	shutdown, err := Init(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer shutdown()

	AddScore(context.Background(), types.Score{BoxName: "box", ServiceName: "svc", TeamNum: 1})

	select {
	case req := <-ch:
		got := req.Header.Get("Authorization")
		if want := "Basic dXNlcjpwYXNz"; got != want {
			t.Errorf("Authorization = %q, want %q", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no request received within timeout")
	}
}

func TestBuildOpts_usernamePassword_encodedCorrectly(t *testing.T) {
	srv, ch := captureServer()
	defer srv.Close()

	cfg := fastBatchCfg(srv.URL)
	cfg.Username = "user"
	cfg.Password = "pass"

	shutdown, err := Init(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer shutdown()

	AddScore(context.Background(), types.Score{BoxName: "box", ServiceName: "svc", TeamNum: 1})

	select {
	case req := <-ch:
		got := req.Header.Get("Authorization")
		if want := "Basic dXNlcjpwYXNz"; got != want { // base64("user:pass")
			t.Errorf("Authorization = %q, want %q", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no request received within timeout")
	}
}

func TestBuildOpts_noAuth_noAuthorizationHeader(t *testing.T) {
	srv, ch := captureServer()
	defer srv.Close()

	shutdown, err := Init(context.Background(), fastBatchCfg(srv.URL))
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer shutdown()

	AddScore(context.Background(), types.Score{BoxName: "box", ServiceName: "svc", TeamNum: 1})

	select {
	case req := <-ch:
		if got := req.Header.Get("Authorization"); got != "" {
			t.Errorf("Authorization = %q, want empty", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no request received within timeout")
	}
}

func TestBuildOpts_streamName_setsHeader(t *testing.T) {
	srv, ch := captureServer()
	defer srv.Close()

	cfg := fastBatchCfg(srv.URL)
	cfg.StreamName = "quincy-test"

	shutdown, err := Init(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer shutdown()

	AddScore(context.Background(), types.Score{BoxName: "box", ServiceName: "svc", TeamNum: 1})

	select {
	case req := <-ch:
		if got := req.Header.Get("stream-name"); got != "quincy-test" {
			t.Errorf("stream-name = %q, want %q", got, "quincy-test")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no request received within timeout")
	}
}
