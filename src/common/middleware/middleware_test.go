package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestPrettyDuration covers all three branches of prettyDuration.
func TestPrettyDuration_microseconds(t *testing.T) {
	s := prettyDuration(500 * time.Nanosecond)
	if s == "" {
		t.Error("expected non-empty string for sub-millisecond duration")
	}
	// should end with µs
	if s[len(s)-3:] != "µs" {
		t.Errorf("expected µs suffix, got %q", s)
	}
}

func TestPrettyDuration_milliseconds(t *testing.T) {
	s := prettyDuration(5 * time.Millisecond)
	if s == "" {
		t.Error("expected non-empty string for millisecond duration")
	}
	if s[len(s)-2:] != "ms" {
		t.Errorf("expected ms suffix, got %q", s)
	}
}

func TestPrettyDuration_seconds(t *testing.T) {
	s := prettyDuration(2 * time.Second)
	if s == "" {
		t.Error("expected non-empty string for second-level duration")
	}
	if s[len(s)-1:] != "s" {
		t.Errorf("expected s suffix, got %q", s)
	}
}

// TestRecovery_catchesPanic verifies Recovery converts panics to 500.
func TestRecovery_catchesPanic(t *testing.T) {
	r := gin.New()
	r.Use(Recovery(false))
	r.GET("/boom", func(c *gin.Context) { panic("test panic") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/boom", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}

func TestRecovery_passesThrough(t *testing.T) {
	r := gin.New()
	r.Use(Recovery(false))
	r.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ok", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

// TestLogging_callsNext verifies Logging passes requests through.
func TestLogging_callsNext(t *testing.T) {
	r := gin.New()
	r.Use(Logging())
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

// TestInsecureCORS_setsHeader verifies the CORS header is set.
func TestInsecureCORS_setsHeader(t *testing.T) {
	r := gin.New()
	r.Use(InsecureCORS())
	r.GET("/data", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/data", nil)
	req.Header.Set("Origin", "http://example.com")
	r.ServeHTTP(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin == "" {
		t.Error("expected Access-Control-Allow-Origin header, got empty")
	}
}
