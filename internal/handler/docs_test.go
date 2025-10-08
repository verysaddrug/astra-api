package handler

import (
	"astra-api/internal/cache"
	mocksgen "astra-api/internal/mocks/gomock"
	"astra-api/internal/model"
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestDocsHandler_List_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().List("u1", gomock.Any()).Return([]model.Document{{ID: "d1", Name: "f", Owner: "u1"}}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=t", nil)

	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_List_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("invalid").Return(model.Session{}, false)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=invalid", nil)

	h.List(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected code 401, got %d", rr.Code)
	}
}

func TestDocsHandler_List_WithLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	userRepo.EXPECT().GetByLogin("otheruser").Return(&model.User{ID: "u2", Login: "otheruser"}, nil)
	docs.EXPECT().List("u2", gomock.Any()).Return([]model.Document{{ID: "d1", Name: "f", Owner: "u2"}}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=t&login=otheruser", nil)

	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_List_WithLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().List("u1", 5).Return([]model.Document{{ID: "d1", Name: "f", Owner: "u1"}}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=t&limit=5", nil)

	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_Upload_JSON_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().Create(gomock.Any()).Return(nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("meta", `{"name":"doc","file":false,"public":false,"mime":"application/json","grants":[]}`)
	_ = mw.WriteField("json", `{"a":1}`)
	_ = mw.WriteField("token", "t")
	_ = mw.Close()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/docs", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	h.Upload(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_Upload_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("invalid").Return(model.Session{}, false)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("meta", `{"name":"doc","file":false,"public":false,"mime":"application/json","grants":[]}`)
	_ = mw.WriteField("token", "invalid")
	_ = mw.Close()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/docs", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	h.Upload(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected code 401, got %d", rr.Code)
	}
}

func TestDocsHandler_Upload_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=t", nil)

	h.Upload(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected code 405, got %d", rr.Code)
	}
}

func TestDocsHandler_Upload_InvalidMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("meta", `invalid json`)
	_ = mw.WriteField("token", "t")
	_ = mw.Close()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/docs", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	h.Upload(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}

func TestDocsHandler_Upload_InvalidJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("meta", `{"name":"doc","file":false,"public":false,"mime":"application/json","grants":[]}`)
	_ = mw.WriteField("json", `invalid json`)
	_ = mw.WriteField("token", "t")
	_ = mw.Close()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/docs", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	h.Upload(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}

func TestDocsHandler_Upload_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().Create(gomock.Any()).Return(errors.New("service error"))

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("meta", `{"name":"doc","file":false,"public":false,"mime":"application/json","grants":[]}`)
	_ = mw.WriteField("json", `{"a":1}`)
	_ = mw.WriteField("token", "t")
	_ = mw.Close()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/docs", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	h.Upload(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected code 500, got %d", rr.Code)
	}
}

func TestDocsHandler_GetByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().GetByID("doc123").Return(&model.Document{
		ID: "doc123", Name: "test.json", Mime: "application/json",
		File: false, Public: false, Owner: "u1", JsonData: []byte(`{"test": "data"}`),
	}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs/doc123?token=t", nil)

	h.GetByID(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_GetByID_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("invalid").Return(model.Session{}, false)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs/doc123?token=invalid", nil)

	h.GetByID(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected code 401, got %d", rr.Code)
	}
}

func TestDocsHandler_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().GetByID("nonexistent").Return(nil, errors.New("not found"))

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs/nonexistent?token=t", nil)

	h.GetByID(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected code 404, got %d", rr.Code)
	}
}

func TestDocsHandler_DeleteByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().Delete("doc123").Return(nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/docs/doc123?token=t", nil)

	h.DeleteByID(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_DeleteByID_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("invalid").Return(model.Session{}, false)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/docs/doc123?token=invalid", nil)

	h.DeleteByID(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected code 401, got %d", rr.Code)
	}
}

func TestDocsHandler_DeleteByID_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs/doc123?token=t", nil)

	h.DeleteByID(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected code 405, got %d", rr.Code)
	}
}

func TestDocsHandler_DeleteByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().Delete("nonexistent").Return(errors.New("not found"))

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/docs/nonexistent?token=t", nil)

	h.DeleteByID(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected code 404, got %d", rr.Code)
	}
}

func TestDocsHandler_List_WithInvalidLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().List("u1", 20).Return([]model.Document{{ID: "d1", Name: "f", Owner: "u1"}}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=t&limit=invalid", nil)

	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_List_WithNegativeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().List("u1", -5).Return([]model.Document{{ID: "d1", Name: "f", Owner: "u1"}}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=t&limit=-5", nil)

	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_List_WithZeroLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().List("u1", 0).Return([]model.Document{{ID: "d1", Name: "f", Owner: "u1"}}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=t&limit=0", nil)

	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_List_WithLargeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().List("u1", 100).Return([]model.Document{{ID: "d1", Name: "f", Owner: "u1"}}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=t&limit=100", nil)

	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_List_WithUnknownLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	userRepo.EXPECT().GetByLogin("unknownuser").Return(nil, errors.New("not found"))

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs?token=t&login=unknownuser", nil)

	h.List(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}

func TestDocsHandler_List_WithHeadMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().List("u1", gomock.Any()).Return([]model.Document{{ID: "d1", Name: "f", Owner: "u1"}}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodHead, "/api/docs?token=t", nil)

	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_GetByID_WithHeadMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)
	docs.EXPECT().GetByID("doc123").Return(&model.Document{
		ID: "doc123", Name: "test.json", Mime: "application/json",
		File: false, Public: false, Owner: "u1", JsonData: []byte(`{"test": "data"}`),
	}, nil)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodHead, "/api/docs/doc123?token=t", nil)

	h.GetByID(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", rr.Code)
	}
}

func TestDocsHandler_GetByID_MissingID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/docs/?token=t", nil)

	h.GetByID(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}

func TestDocsHandler_GetByID_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/docs/doc123?token=t", nil)

	h.GetByID(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected code 405, got %d", rr.Code)
	}
}

func TestDocsHandler_DeleteByID_MissingID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docs := mocksgen.NewMockDocsServiceInterface(ctrl)
	sess := mocksgen.NewMockSessionServiceInterface(ctrl)
	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)

	sess.EXPECT().Validate("t").Return(model.Session{UserID: "u1"}, true)

	c := cache.NewCache(0)
	h := NewDocsHandler(docs, c, sess, userRepo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/docs/?token=t", nil)

	h.DeleteByID(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code 400, got %d", rr.Code)
	}
}
