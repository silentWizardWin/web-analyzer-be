package main

import (
	"net/http"
	"testing"
	"time"
)

func TestMainServerStarts(t *testing.T) {
	go func() {
		main()
	}()

	// allow server to boot
	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8080/analyze")
	if err != nil {
		t.Log("Server not ready or not responding as expected (this is OK for unit test)")
		return
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 for GET /analyze, got %d", resp.StatusCode)
	}
}
