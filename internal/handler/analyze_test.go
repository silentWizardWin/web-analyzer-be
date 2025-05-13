package handler

import (
	"bytes"
	"strings"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalyzeHandlerInvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/analyze", nil)
	w := httptest.NewRecorder()

	AnalyzeHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", w.Code)
	}
}

func TestAnalyzeHandlerBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer([]byte("invalid")))
	w := httptest.NewRecorder()

	AnalyzeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestAnalyzeHandlerValidRequest(t *testing.T) {
	// Use test server
	html := `<html><head><title>Test</title></head><body><h1>H</h1></body></html>`
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	defer testServer.Close()

	body, _ := json.Marshal(map[string]string{"url": testServer.URL})
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	AnalyzeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("Invalid JSON response: %v", err)
	}
	if result["title"] != "Test" {
		t.Errorf("Expected title 'Test', got %v", result["title"])
	}
}

func TestAnalyzeHandlerInvalidJSON(t *testing.T) {
	body := []byte(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	AnalyzeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", w.Code)
	}
}

func TestAnalyzeHandlerMissingURL(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"url": ""})
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	AnalyzeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty URL, got %d", w.Code)
	}

	expected := `{"error":"Invalid JSON: missing or malformed 'url' field","status":400}`
	if strings.TrimSpace(w.Body.String()) != expected {
		t.Errorf("Expected response body:\n%s\ngot:\n%s", expected, w.Body.String())
	}
}

func TestAnalyzeHandlerBadURL(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"url": "http://invalid.localhost.test"})
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	AnalyzeHandler(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("Expected 502 for unreachable URL, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if resp["status"] != float64(http.StatusBadGateway) {
		t.Errorf("Expected status 502 in response body, got %v", resp["status"])
	}

	if _, ok := resp["error"]; !ok {
		t.Error("Expected 'error' field in response body")
	}
}
