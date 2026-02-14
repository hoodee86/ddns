package main

import (
"net/http"
"net/http/httptest"
"testing"
)

func TestGetClientIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	clientIP := getClientIP(req)
	if clientIP != "192.168.1.100" {
		t.Errorf("expected 192.168.1.100, got %s", clientIP)
	}
}

func TestHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "203.0.113.42:12345"
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
