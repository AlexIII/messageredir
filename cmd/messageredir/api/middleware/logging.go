package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger is a middleware function that logs HTTP requests.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log the request details
		log.Printf(
			"%s %s took %s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
