package main

import (
	"fmt"
	"net/http"
)

// recoverPanic is a middleware that catches any panics in downstream handlers
// and returns a 500 JSON error response instead of silently dropping the connection.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				// Close the connection and send a 500 JSON response
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
