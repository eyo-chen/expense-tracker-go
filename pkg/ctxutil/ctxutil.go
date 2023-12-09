package ctxutil

import (
	"context"
	"net/http"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type contextKey string

const contextKeyUser = contextKey("user")

// Set stores the user in the request context
func SetUser(r *http.Request, user *domain.User) *http.Request {
	ctx := context.WithValue(r.Context(), contextKeyUser, user)
	return r.WithContext(ctx)
}

// Get retrieves the user from the request context
func GetUser(r *http.Request) *domain.User {
	user, ok := r.Context().Value(contextKeyUser).(*domain.User)
	if !ok {
		logger.Panic("missing user value in request context")
	}

	return user
}
