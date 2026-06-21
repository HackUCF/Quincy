package config

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func configuredSinks() Sinks {
	return Sinks{PGConfig: PGConfig{Host: "localhost", Port: 5432, Username: "u", Password: "p", Database: "d"}}
}

func callHandler(h gin.HandlerFunc) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	h(c)
	return w
}

func TestSinks_DBEnabled_empty(t *testing.T) {
	s := Sinks{}
	if s.DBEnabled() {
		t.Error("DBEnabled() should be false when no PGConfig is set")
	}
}

func TestSinks_DBEnabled_configured(t *testing.T) {
	s := configuredSinks()
	if !s.DBEnabled() {
		t.Error("DBEnabled() should be true when PGConfig is set")
	}
}

func TestSinks_DBOr501_withDB_returnsEndpoint(t *testing.T) {
	s := configuredSinks()
	sentinel := gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	w := callHandler(s.DBOr501(sentinel))
	if w.Code != http.StatusOK {
		t.Errorf("DBOr501 with configured DB: got %d, want 200", w.Code)
	}
}

func TestSinks_DBOr501_noDB_returns501(t *testing.T) {
	s := Sinks{}
	sentinel := gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	w := callHandler(s.DBOr501(sentinel))
	if w.Code != http.StatusNotImplemented {
		t.Errorf("DBOr501 without DB: got %d, want 501", w.Code)
	}
}

func TestSinks_DBOrNOP_withDB_callsHandler(t *testing.T) {
	s := configuredSinks()
	called := false
	handler := gin.HandlerFunc(func(c *gin.Context) { called = true })
	callHandler(s.DBOrNOP(handler))
	if !called {
		t.Error("DBOrNOP with configured DB should call the handler")
	}
}

func TestSinks_DBOrNOP_noDB_isNOP(t *testing.T) {
	s := Sinks{}
	called := false
	handler := gin.HandlerFunc(func(c *gin.Context) { called = true })
	callHandler(s.DBOrNOP(handler))
	if called {
		t.Error("DBOrNOP without DB should not call the handler")
	}
}

func TestSinks_OTelEnabled_empty(t *testing.T) {
	s := Sinks{}
	if s.OTelEnabled() {
		t.Error("OTelEnabled() should be false when no endpoint is set")
	}
}

func TestSinks_OTelEnabled_configured(t *testing.T) {
	s := Sinks{OTelConfig: OTelConfig{Endpoint: "http://localhost:4318"}}
	if !s.OTelEnabled() {
		t.Error("OTelEnabled() should be true when endpoint is set")
	}
}
