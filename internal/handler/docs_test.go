package handler

import (
	"astra-api/internal/cache"
	mocksgen "astra-api/internal/mocks/gomock"
	"astra-api/internal/model"
	"bytes"
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
