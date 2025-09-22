package middleware

import (
	"astra-api/internal/interfaces"
	"astra-api/internal/model"
	"context"
	"log"
	"net/http"
	"time"
)

// ContextKey represents a key used for context values
type ContextKey string

const (
	// UserContextKey is the key for storing user in context
	UserContextKey ContextKey = "user"
	// SessionContextKey is the key for storing session in context  
	SessionContextKey ContextKey = "session"
)

// AuthMiddleware provides authentication middleware
type AuthMiddleware struct {
	sessionService interfaces.SessionServiceInterface
	userRepo       interfaces.UserRepositoryInterface
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(sessionService interfaces.SessionServiceInterface, userRepo interfaces.UserRepositoryInterface) *AuthMiddleware {
	return &AuthMiddleware{
		sessionService: sessionService,
		userRepo:       userRepo,
	}
}

// RequireAuth is middleware that requires authentication
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		if token == "" {
			http.Error(w, `{"error":"missing authentication token"}`, http.StatusUnauthorized)
			return
		}

		session, exists := m.sessionService.Get(token)
		if !exists {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		user, err := m.userRepo.GetByID(session.UserID)
		if err != nil {
			http.Error(w, `{"error":"user not found"}`, http.StatusUnauthorized)
			return
		}

		// Add user and session to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, SessionContextKey, session)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserFromContext extracts user from request context
func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*model.User)
	return user, ok
}

// GetSessionFromContext extracts session from request context
func GetSessionFromContext(ctx context.Context) (*model.Session, bool) {
	session, ok := ctx.Value(SessionContextKey).(*model.Session)
	return session, ok
}

// getToken extracts token from request (query param, header, or form)
func getToken(r *http.Request) string {
	if token := r.URL.Query().Get("token"); token != "" {
		return token
	}
	if token := r.Header.Get("Authorization"); token != "" {
		return token
	}
	if err := r.ParseMultipartForm(32 << 20); err == nil {
		if token := r.FormValue("token"); token != "" {
			return token
		}
	}
	return ""
}