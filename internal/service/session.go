package service

import (
	"astra-api/internal/model"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SessionService struct {
	mu       sync.RWMutex
	sessions map[string]model.Session
}

func NewSessionService() *SessionService {
	return &SessionService{sessions: make(map[string]model.Session)}
}

func (s *SessionService) Create(userID, login string) string {
	token := uuid.New().String()
	s.mu.Lock()
	s.sessions[token] = model.Session{Token: token, UserID: userID, Login: login, Created: time.Now()}
	s.mu.Unlock()
	return token
}

func (s *SessionService) Get(token string) (*model.Session, bool) {
	s.mu.RLock()
	sess, ok := s.sessions[token]
	s.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return &sess, true
}

func (s *SessionService) Validate(token string) (model.Session, bool) {
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
