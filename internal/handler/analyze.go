package handler

import (
	"encoding/json"
	"io"
	"fmt"
	"net/http"
	"strings"
	"web-analyzer-be/internal/model"
	"web-analyzer-be/internal/service"
)

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed","status":405}`, http.StatusMethodNotAllowed)
		return
	}

	var req model.AnalyzeRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read request","status":400}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil || req.URL == "" {
		http.Error(w, `{"error":"Invalid JSON: missing or malformed 'url' field","status":400}`, http.StatusBadRequest)
		return
	}

	resp, err := service.AnalyzeURL(req.URL)
	if err != nil {
		status := http.StatusInternalServerError
		msg := err.Error()

		switch {
		case strings.Contains(msg, "no such host"):
			status = http.StatusBadGateway
		case strings.Contains(msg, "timeout"):
			status = http.StatusGatewayTimeout
		case strings.Contains(msg, "HTTP status"):
			status = http.StatusBadGateway
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  fmt.Sprintf("Failed to analyze page: %s", msg),
			"status": status,
		})
		return
	}

	json.NewEncoder(w).Encode(resp)
}
