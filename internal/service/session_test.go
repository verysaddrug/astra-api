package service

import (
	"testing"
	"time"
)

func TestSessionService_Create(t *testing.T) {
	sessionService := NewSessionService()

	token := sessionService.Create("user123", "testuser")

	if token == "" {
		t.Fatal("expected token to be generated")
	}

	// Verify session was created
	session, ok := sessionService.Get(token)
	if !ok {
		t.Fatal("expected session to exist")
	}
	if session.UserID != "user123" {
		t.Fatalf("expected UserID 'user123', got %s", session.UserID)
	}
	if session.Login != "testuser" {
		t.Fatalf("expected Login 'testuser', got %s", session.Login)
	}
	if session.Token != token {
		t.Fatalf("expected Token '%s', got %s", token, session.Token)
	}
}

func TestSessionService_Get_Existing(t *testing.T) {
	sessionService := NewSessionService()

	token := sessionService.Create("user123", "testuser")

	session, ok := sessionService.Get(token)

	if !ok {
		t.Fatal("expected session to exist")
	}
	if session.UserID != "user123" {
		t.Fatalf("expected UserID 'user123', got %s", session.UserID)
	}
}

func TestSessionService_Get_NonExistent(t *testing.T) {
	sessionService := NewSessionService()

	session, ok := sessionService.Get("nonexistent")

	if ok {
		t.Fatal("expected session to not exist")
	}
	if session != nil {
		t.Fatal("expected session to be nil")
	}
}

func TestSessionService_Validate_Existing(t *testing.T) {
	sessionService := NewSessionService()

	token := sessionService.Create("user123", "testuser")

	session, ok := sessionService.Validate(token)

	if !ok {
		t.Fatal("expected session to exist")
	}
	if session.UserID != "user123" {
		t.Fatalf("expected UserID 'user123', got %s", session.UserID)
	}
}

func TestSessionService_Validate_NonExistent(t *testing.T) {
	sessionService := NewSessionService()

	session, ok := sessionService.Validate("nonexistent")

	if ok {
		t.Fatal("expected session to not exist")
	}
	if session.UserID != "" {
		t.Fatal("expected empty UserID")
	}
}

func TestSessionService_Delete_Existing(t *testing.T) {
	sessionService := NewSessionService()

	token := sessionService.Create("user123", "testuser")

	deleted := sessionService.Delete(token)

	if !deleted {
		t.Fatal("expected session to be deleted")
	}

	// Verify session was deleted
	_, ok := sessionService.Get(token)
	if ok {
		t.Fatal("expected session to be deleted")
	}
}

func TestSessionService_Delete_NonExistent(t *testing.T) {
	sessionService := NewSessionService()

	deleted := sessionService.Delete("nonexistent")

	if deleted {
		t.Fatal("expected session to not exist")
	}
}

func TestSessionService_ConcurrentAccess(t *testing.T) {
	sessionService := NewSessionService()

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			token := sessionService.Create("user"+string(rune(i)), "testuser")
			sessionService.Validate(token)
			sessionService.Delete(token)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestSessionService_CreatedAt(t *testing.T) {
	sessionService := NewSessionService()

	before := time.Now()
	token := sessionService.Create("user123", "testuser")
	after := time.Now()

	session, ok := sessionService.Get(token)
	if !ok {
		t.Fatal("expected session to exist")
	}

	if session.Created.Before(before) || session.Created.After(after) {
		t.Fatal("expected CreatedAt to be within expected time range")
	}
}
