package http

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/anubhav047/goboard/internal/db"
)

// contextKey is a custom type to avoid key collision in context.
type contextKey string

const userContextKey = contextKey("user")

// Middleware struct holds dependencies for middleware
type Middleware struct {
	sm      *scs.SessionManager
	queries *db.Queries
}

// NewMiddleware creates a new Middleware struct.
func NewMiddleware(sm *scs.SessionManager, queries *db.Queries) *Middleware {
	return &Middleware{
		sm:      sm,
		queries: queries,
	}
}

// RequireAuth is the middleware that protects routes
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		userID := m.sm.GetInt32(r.Context(), "authenticatedUserID")
		if userID == 0 {
			WriteError(w, http.StatusUnauthorized, "You must be logged in to access this resourse")
			return
		}

		// Fetch the user from the database
		user, err := m.queries.GetUserByID(r.Context(), userID)
		if err != nil {
			// This could happen if the user was deleted after the session was created.
			WriteError(w, http.StatusUnauthorized, "Invalid authentication token")
			return
		}

		// Add the user to request context
		ctx := context.WithValue(r.Context(), userContextKey, user)

		// Call the next handler in the chain, using the new context.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
