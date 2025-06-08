package main

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingMiddleware logs the time and URL of each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Middleware type to allow chaining multiple middleware functions
type Middleware func(http.Handler) http.Handler

// ChainMiddleware chains multiple middleware functions together
func ChainMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	// Apply multiple middlewares using ChainMiddleware
	middlewares := []Middleware{
		LoggingMiddleware,
		// Add more middleware here
	}
	wrappedMux := ChainMiddleware(mux, middlewares...)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
