package service

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token   string
	UserID  string
	Login   string
	Created time.Time
}

type SessionService struct {
	mu       sync.RWMutex
	sessions map[string]Session
}

func NewSessionService() *SessionService {
	return &SessionService{sessions: make(map[string]Session)}
}

func (s *SessionService) Create(userID, login string) string {
	token := uuid.New().String()
	s.mu.Lock()
	s.sessions[token] = Session{Token: token, UserID: userID, Login: login, Created: time.Now()}
	s.mu.Unlock()
	return token
}

func (s *SessionService) Validate(token string) (Session, bool) {
	s.mu.RLock()
	sess, ok := s.sessions[token]
	s.mu.RUnlock()
	return sess, ok
}

func (s *SessionService) Delete(token string) bool {
	s.mu.Lock()
	_, ok := s.sessions[token]
	if ok {
		delete(s.sessions, token)
	}
	s.mu.Unlock()
	return ok
}
