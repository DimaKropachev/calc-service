package transport

import (
	"log"
	"net/http"
	"time"
)

// Middleware для отловки паники
func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] recovered from panic: %v, Method: %v, Path: %v\n\n", r.Method, r.URL.Path, err)
				SendErrorJSON(w, http.StatusText(http.StatusInternalServerError))
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Middleware для логирования данных о HTTP запросе
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("Request:\nMethod: %v\nPath: %v\nRequest time: %v\n\n", r.Method, r.URL.Path, time.Since(start))
	})
}
