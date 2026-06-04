package main

import (
	"context"
	"net/http"

	"github.com/ebenamoafo2/scalable-go-api/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

// ContextSetUser returns a new copy of the request with the
// provided User struct added to the context
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// ContextGetUser retrieves the User struct from the request context
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in the request context")
	}

	return user
}
