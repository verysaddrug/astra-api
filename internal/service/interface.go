package service

import "astra-api/internal/model"

// AuthServiceInterface описывает контракт сервиса аутентификации
type AuthServiceInterface interface {
	Register(login, password, adminToken string) (*model.User, error)
	Authenticate(login, password string) (*model.User, error)
}

// DocsServiceInterface описывает контракт сервиса документов
type DocsServiceInterface interface {
	Create(doc *model.Document) error
	List(owner string, limit int) ([]model.Document, error)
	GetByID(id string) (*model.Document, error)
	Delete(id string) error
}

// SessionServiceInterface описывает контракт сервиса сессий
type SessionServiceInterface interface {
	Create(userID, login string) string
	Get(token string) (*model.Session, bool)
	Validate(token string) (model.Session, bool)
	Delete(token string) bool
}
