package handler

import (
	mocksgen "astra-api/internal/mocks/gomock"
	"astra-api/internal/model"
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
