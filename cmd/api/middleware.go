package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
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

func (app *application) rateLimit(next http.Handler) http.Handler {
	//Define a custom type for each client
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	//Periodically check for expired clients and remove them from the map
	go func() {
		time.Sleep(time.Minute)
		mu.Lock()

		//delete any clients that have been inactive for more than 3 minutes
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.config.limiter.enabled {
			ip := realip.FromRequest(r)

			//Lock the mutex to prevent concurrent access to the client's map
			mu.Lock()

			//check if the ip address exists on the map, if not, create one
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
			}
			//Update the last seen time for the client
			clients[ip].lastSeen = time.Now()

			//Unlock the mutex and check if the client is allowed to make a request
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
		}

		//If the client is allowed to make a request, proceed to the next handler
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
