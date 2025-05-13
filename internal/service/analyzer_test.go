package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalyzeURL(t *testing.T) {
	html := `
	<!DOCTYPE html>
	<html><head><title>Test</title></head>
	<body><h1>Hello</h1><form><input type="password"/></form></body></html>
	`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	resp, err := AnalyzeURL(server.URL)
	if err != nil {
		t.Fatalf("AnalyzeURL failed: %v", err)
	}

	if resp.Title != "Test" {
		t.Errorf("Expected title 'Test', got '%s'", resp.Title)
	}
	if !resp.LoginFormExists {
		t.Error("Expected login form to be detected")
	}
}
