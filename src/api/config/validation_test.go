package config

import "testing"

func TestHTTPSpec_validate_valid(t *testing.T) {
	h := HTTPSpec{Host: "0.0.0.0", Port: 8888}
	if err := h.validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestHTTPSpec_validate_portNegative(t *testing.T) {
	h := HTTPSpec{Host: "0.0.0.0", Port: -1}
	if err := h.validate(); err == nil {
		t.Fatal("expected error for port -1")
	}
}

func TestHTTPSpec_validate_portTooHigh(t *testing.T) {
	h := HTTPSpec{Host: "0.0.0.0", Port: 99999}
	if err := h.validate(); err == nil {
		t.Fatal("expected error for port 99999")
	}
}

func TestHTTPSpec_validate_badHost(t *testing.T) {
	h := HTTPSpec{Host: "not-an-ip", Port: 8888}
	if err := h.validate(); err == nil {
		t.Fatal("expected error for non-IP host")
	}
}

func TestHTTPSpec_validate_ipv6(t *testing.T) {
	h := HTTPSpec{Host: "::1", Port: 8888}
	if err := h.validate(); err != nil {
		t.Fatalf("expected no error for ipv6, got: %v", err)
	}
}

func TestHTTPSpec_validate_portZero(t *testing.T) {
	h := HTTPSpec{Host: "127.0.0.1", Port: 0}
	if err := h.validate(); err != nil {
		t.Fatalf("expected no error for port 0 (valid TCP port), got: %v", err)
	}
}
