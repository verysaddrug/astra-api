package handler

import (
	mocksgen "astra-api/internal/mocks/gomock"
	"astra-api/internal/model"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Register_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	auth.EXPECT().Register("a", "b", "adm").Return(&model.User{Login: "a", ID: "u1"}, nil)

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(`{"login":"a","pswd":"b","token":"adm"}`))
	req.Header.Set("Content-Type", "application/json")

	h.Register(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestAuthHandler_Auth_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	auth.EXPECT().Authenticate("a", "b").Return(&model.User{Login: "a", ID: "u1"}, nil)
	sess.EXPECT().Create("u1", "a").Return("tok")

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(`{"login":"a","pswd":"b"}`))
	req.Header.Set("Content-Type", "application/json")

	h.Auth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestAuthHandler_Logout_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	sess.EXPECT().Delete("tok").Return(true)

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/auth/tok", nil)

	h.Logout(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestAuthHandler_Register_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	auth.EXPECT().Register("a", "b", "wrong").Return(nil, errors.New("invalid admin token"))

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(`{"login":"a","pswd":"b","token":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")

	h.Register(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")

	h.Register(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}

func TestAuthHandler_Register_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/register", nil)

	h.Register(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected code 405, got %d", rr.Code)
	}
}

func TestAuthHandler_Auth_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	auth.EXPECT().Authenticate("a", "wrong").Return(nil, errors.New("invalid password"))

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(`{"login":"a","pswd":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")

	h.Auth(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected code 401, got %d", rr.Code)
	}
}

func TestAuthHandler_Auth_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")

	h.Auth(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}

func TestAuthHandler_Auth_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/auth", nil)

	h.Auth(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected code 405, got %d", rr.Code)
	}
}

func TestAuthHandler_Logout_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	sess.EXPECT().Delete("invalid").Return(false)

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/auth/invalid", nil)

	h.Logout(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}

func TestAuthHandler_Logout_MissingToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	// No expectations for Delete since the token is missing

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/auth/", nil)

	h.Logout(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}

func TestAuthHandler_Logout_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocksgen.NewMockAuthServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)

	h := NewAuthHandler(auth, sess)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/auth/tok", nil)

	h.Logout(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected code 405, got %d", rr.Code)
	}
}
