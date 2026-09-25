// backend/tests/api_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpointExists(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	// тело должно содержать поле "status"
	if !contains(rec.Body.String(), `"status"`) {
		t.Fatalf("body %q missing 'status'", rec.Body.String())
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})())
}
