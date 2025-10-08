package model

import (
	"testing"
	"time"

	"github.com/lib/pq"
)

func TestAPIError(t *testing.T) {
	error := APIError{
		Code: 400,
		Text: "Bad Request",
	}

	if error.Code != 400 {
		t.Fatalf("expected Code 400, got %d", error.Code)
	}
	if error.Text != "Bad Request" {
		t.Fatalf("expected Text 'Bad Request', got %s", error.Text)
	}
}

func TestAPIResponse_WithError(t *testing.T) {
	response := APIResponse{
		Error: &APIError{
			Code: 404,
			Text: "Not Found",
		},
	}

	if response.Error == nil {
		t.Fatal("expected Error to be set")
	}
	if response.Error.Code != 404 {
		t.Fatalf("expected Error.Code 404, got %d", response.Error.Code)
	}
	if response.Error.Text != "Not Found" {
		t.Fatalf("expected Error.Text 'Not Found', got %s", response.Error.Text)
	}
}

func TestAPIResponse_WithResponse(t *testing.T) {
	response := APIResponse{
		Response: map[string]string{
			"message": "success",
		},
	}

	if response.Response == nil {
		t.Fatal("expected Response to be set")
	}
	if response.Error != nil {
		t.Fatal("expected Error to be nil")
	}
}

func TestAPIResponse_WithData(t *testing.T) {
	response := APIResponse{
		Data: []string{"item1", "item2"},
	}

	if response.Data == nil {
		t.Fatal("expected Data to be set")
	}
	if response.Error != nil {
		t.Fatal("expected Error to be nil")
	}
}

func TestSession(t *testing.T) {
	session := Session{
		Token:   "token123",
		UserID:  "user123",
		Login:   "testuser",
		Created: time.Now(),
	}

	if session.Token != "token123" {
		t.Fatalf("expected Token 'token123', got %s", session.Token)
	}
	if session.UserID != "user123" {
		t.Fatalf("expected UserID 'user123', got %s", session.UserID)
	}
	if session.Login != "testuser" {
		t.Fatalf("expected Login 'testuser', got %s", session.Login)
	}
	if session.Created.IsZero() {
		t.Fatal("expected Created to be set")
	}
}

func TestDocument(t *testing.T) {
	grants := pq.StringArray{"read", "write"}
	jsonData := []byte(`{"test": "data"}`)

	doc := Document{
		ID:        "doc123",
		Name:      "test.json",
		Mime:      "application/json",
		File:      false,
		Public:    false,
		Owner:     "user123",
		CreatedAt: time.Now(),
		Grants:    grants,
		JsonData:  jsonData,
	}

	if doc.ID != "doc123" {
		t.Fatalf("expected ID 'doc123', got %s", doc.ID)
	}
	if doc.Name != "test.json" {
		t.Fatalf("expected Name 'test.json', got %s", doc.Name)
	}
	if doc.Mime != "application/json" {
		t.Fatalf("expected Mime 'application/json', got %s", doc.Mime)
	}
	if doc.File {
		t.Fatal("expected File to be false")
	}
	if doc.Public {
		t.Fatal("expected Public to be false")
	}
	if doc.Owner != "user123" {
		t.Fatalf("expected Owner 'user123', got %s", doc.Owner)
	}
	if doc.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
	if len(doc.Grants) != 2 {
		t.Fatalf("expected 2 grants, got %d", len(doc.Grants))
	}
	if doc.Grants[0] != "read" {
		t.Fatalf("expected first grant 'read', got %s", doc.Grants[0])
	}
	if doc.Grants[1] != "write" {
		t.Fatalf("expected second grant 'write', got %s", doc.Grants[1])
	}
	if len(doc.JsonData) == 0 {
		t.Fatal("expected JsonData to be set")
	}
}

func TestDocument_EmptyGrants(t *testing.T) {
	doc := Document{
		ID:     "doc123",
		Name:   "test.json",
		Grants: pq.StringArray{},
	}

	if len(doc.Grants) != 0 {
		t.Fatalf("expected 0 grants, got %d", len(doc.Grants))
	}
}

func TestDocument_EmptyJsonData(t *testing.T) {
	doc := Document{
		ID:       "doc123",
		Name:     "test.json",
		JsonData: []byte{},
	}

	if len(doc.JsonData) != 0 {
		t.Fatalf("expected empty JsonData, got %d bytes", len(doc.JsonData))
	}
}

func TestDocument_FileType(t *testing.T) {
	doc := Document{
		ID:   "doc123",
		Name: "test.pdf",
		Mime: "application/pdf",
		File: true,
	}

	if !doc.File {
		t.Fatal("expected File to be true")
	}
	if doc.Mime != "application/pdf" {
		t.Fatalf("expected Mime 'application/pdf', got %s", doc.Mime)
	}
}

func TestDocument_Public(t *testing.T) {
	doc := Document{
		ID:     "doc123",
		Name:   "public.json",
		Public: true,
	}

	if !doc.Public {
		t.Fatal("expected Public to be true")
	}
}
