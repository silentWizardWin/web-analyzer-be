package main

import (
	"log"
	"net/http"

	"web-analyzer-be/internal/handler"
)

// enable CORS for frontend
func enableCORS(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            return
        }
        h.ServeHTTP(w, r)
    })
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/analyze", handler.AnalyzeHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(router)))
}
