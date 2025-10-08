package service

import (
	mocksgen "astra-api/internal/mocks/gomock"
	"astra-api/internal/model"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestDocsService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docRepo := mocksgen.NewMockDocumentRepositoryInterface(ctrl)
	docsService := NewDocsService(docRepo)

	doc := &model.Document{
		Name:     "test.json",
		Mime:     "application/json",
		File:     false,
		Public:   false,
		Owner:    "user123",
		Grants:   []string{},
		JsonData: []byte(`{"test": "data"}`),
	}

	docRepo.EXPECT().Create(gomock.Any()).Do(func(doc *model.Document) {
		if doc.ID == "" {
			t.Fatal("expected ID to be set")
		}
		if doc.CreatedAt.IsZero() {
			t.Fatal("expected CreatedAt to be set")
		}
	}).Return(nil)

	err := docsService.Create(doc)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDocsService_Create_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docRepo := mocksgen.NewMockDocumentRepositoryInterface(ctrl)
	docsService := NewDocsService(docRepo)

	doc := &model.Document{
		Name:     "test.json",
		Mime:     "application/json",
		File:     false,
		Public:   false,
		Owner:    "user123",
		Grants:   []string{},
		JsonData: []byte(`{"test": "data"}`),
	}

	docRepo.EXPECT().Create(gomock.Any()).Return(errors.New("database error"))

	err := docsService.Create(doc)

	if err == nil {
		t.Fatal("expected error from repository")
	}
	if err.Error() != "database error" {
		t.Fatalf("expected 'database error', got %s", err.Error())
	}
}

func TestDocsService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docRepo := mocksgen.NewMockDocumentRepositoryInterface(ctrl)
	docsService := NewDocsService(docRepo)

	expectedDocs := []model.Document{
		{
			ID:        uuid.New().String(),
			Name:      "doc1.json",
			Mime:      "application/json",
			File:      false,
			Public:    false,
			Owner:     "user123",
			CreatedAt: time.Now(),
			Grants:    []string{},
			JsonData:  []byte(`{"test": "data1"}`),
		},
		{
			ID:        uuid.New().String(),
			Name:      "doc2.json",
			Mime:      "application/json",
			File:      false,
			Public:    false,
			Owner:     "user123",
			CreatedAt: time.Now(),
			Grants:    []string{},
			JsonData:  []byte(`{"test": "data2"}`),
		},
	}

	docRepo.EXPECT().List("user123", 10).Return(expectedDocs, nil)

	docs, err := docsService.List("user123", 10)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(docs) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(docs))
	}
	if docs[0].Name != "doc1.json" {
		t.Fatalf("expected first doc name 'doc1.json', got %s", docs[0].Name)
	}
}

func TestDocsService_List_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docRepo := mocksgen.NewMockDocumentRepositoryInterface(ctrl)
	docsService := NewDocsService(docRepo)

	docRepo.EXPECT().List("user123", 10).Return(nil, errors.New("database error"))

	docs, err := docsService.List("user123", 10)

	if err == nil {
		t.Fatal("expected error from repository")
	}
	if docs != nil {
		t.Fatal("expected docs to be nil")
	}
	if err.Error() != "database error" {
		t.Fatalf("expected 'database error', got %s", err.Error())
	}
}

func TestDocsService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docRepo := mocksgen.NewMockDocumentRepositoryInterface(ctrl)
	docsService := NewDocsService(docRepo)

	expectedDoc := &model.Document{
		ID:        "doc123",
		Name:      "test.json",
		Mime:      "application/json",
		File:      false,
		Public:    false,
		Owner:     "user123",
		CreatedAt: time.Now(),
		Grants:    []string{},
		JsonData:  []byte(`{"test": "data"}`),
	}

	docRepo.EXPECT().GetByID("doc123").Return(expectedDoc, nil)

	doc, err := docsService.GetByID("doc123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if doc.ID != "doc123" {
		t.Fatalf("expected doc ID 'doc123', got %s", doc.ID)
	}
	if doc.Name != "test.json" {
		t.Fatalf("expected doc name 'test.json', got %s", doc.Name)
	}
}

func TestDocsService_GetByID_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docRepo := mocksgen.NewMockDocumentRepositoryInterface(ctrl)
	docsService := NewDocsService(docRepo)

	docRepo.EXPECT().GetByID("nonexistent").Return(nil, errors.New("not found"))

	doc, err := docsService.GetByID("nonexistent")

	if err == nil {
		t.Fatal("expected error from repository")
	}
	if doc != nil {
		t.Fatal("expected doc to be nil")
	}
	if err.Error() != "not found" {
		t.Fatalf("expected 'not found', got %s", err.Error())
	}
}

func TestDocsService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docRepo := mocksgen.NewMockDocumentRepositoryInterface(ctrl)
	docsService := NewDocsService(docRepo)

	docRepo.EXPECT().Delete("doc123").Return(nil)

	err := docsService.Delete("doc123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDocsService_Delete_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	docRepo := mocksgen.NewMockDocumentRepositoryInterface(ctrl)
	docsService := NewDocsService(docRepo)

	docRepo.EXPECT().Delete("nonexistent").Return(errors.New("not found"))

	err := docsService.Delete("nonexistent")

	if err == nil {
		t.Fatal("expected error from repository")
	}
	if err.Error() != "not found" {
		t.Fatalf("expected 'not found', got %s", err.Error())
	}
}
