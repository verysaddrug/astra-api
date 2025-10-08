package middleware

import (
	mocksgen "astra-api/internal/mocks/gomock"
	"astra-api/internal/model"
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware_RequireAuth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionService := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authMiddleware := NewAuthMiddleware(sessionService, userRepo)

	user := &model.User{
		ID:    uuid.New().String(),
		Login: "testuser",
	}

	session := &model.Session{
		Token:  "validtoken",
		UserID: user.ID,
		Login:  user.Login,
	}

	sessionService.EXPECT().Get("validtoken").Return(session, true)
	userRepo.EXPECT().GetByID(user.ID).Return(user, nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is in context
		ctxUser := r.Context().Value(UserContextKey)
		if ctxUser == nil {
			t.Fatal("expected user in context")
		}
		// Check if session is in context
		ctxSession := r.Context().Value(SessionContextKey)
		if ctxSession == nil {
			t.Fatal("expected session in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := authMiddleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/test?token=validtoken", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestAuthMiddleware_RequireAuth_MissingToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionService := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authMiddleware := NewAuthMiddleware(sessionService, userRepo)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	handler := authMiddleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_RequireAuth_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionService := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authMiddleware := NewAuthMiddleware(sessionService, userRepo)

	sessionService.EXPECT().Get("invalidtoken").Return(nil, false)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	handler := authMiddleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/test?token=invalidtoken", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_RequireAuth_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionService := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authMiddleware := NewAuthMiddleware(sessionService, userRepo)

	session := &model.Session{
		Token:  "validtoken",
		UserID: "nonexistent",
		Login:  "testuser",
	}

	sessionService.EXPECT().Get("validtoken").Return(session, true)
	userRepo.EXPECT().GetByID("nonexistent").Return(nil, &model.UserNotFoundError{})

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	handler := authMiddleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/test?token=validtoken", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_RequireAuth_TokenFromHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionService := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authMiddleware := NewAuthMiddleware(sessionService, userRepo)

	user := &model.User{
		ID:    uuid.New().String(),
		Login: "testuser",
	}

	session := &model.Session{
		Token:  "headertoken",
		UserID: user.ID,
		Login:  user.Login,
	}

	sessionService.EXPECT().Get("headertoken").Return(session, true)
	userRepo.EXPECT().GetByID(user.ID).Return(user, nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := authMiddleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "headertoken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestAuthMiddleware_RequireAuth_TokenFromForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionService := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authMiddleware := NewAuthMiddleware(sessionService, userRepo)

	user := &model.User{
		ID:    uuid.New().String(),
		Login: "testuser",
	}

	session := &model.Session{
		Token:  "formtoken",
		UserID: user.ID,
		Login:  user.Login,
	}

	sessionService.EXPECT().Get("formtoken").Return(session, true)
	userRepo.EXPECT().GetByID(user.ID).Return(user, nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := authMiddleware.RequireAuth(next)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("token", "formtoken")
	_ = mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/test", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := LoggingMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	wrapper := &responseWriter{ResponseWriter: rr, statusCode: http.StatusOK}

	wrapper.WriteHeader(http.StatusNotFound)

	if wrapper.statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", wrapper.statusCode)
	}
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected recorder status 404, got %d", rr.Code)
	}
}

func TestChainMiddleware(t *testing.T) {
	var calls []string

	middleware1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "middleware1")
			next.ServeHTTP(w, r)
		}
	}

	middleware2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "middleware2")
			next.ServeHTTP(w, r)
		}
	}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "final")
		w.WriteHeader(http.StatusOK)
	})

	handler := ChainMiddleware(middleware1, middleware2)(final)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	expectedCalls := []string{"middleware1", "middleware2", "final"}
	if len(calls) != len(expectedCalls) {
		t.Fatalf("expected %d calls, got %d", len(expectedCalls), len(calls))
	}
	for i, expected := range expectedCalls {
		if calls[i] != expected {
			t.Fatalf("expected call %d to be '%s', got '%s'", i, expected, calls[i])
		}
	}
}
