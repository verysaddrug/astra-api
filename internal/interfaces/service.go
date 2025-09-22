package interfaces

import "astra-api/internal/model"

// AuthServiceInterface defines the contract for authentication service
type AuthServiceInterface interface {
	Register(login, password, adminToken string) (*model.User, error)
	Authenticate(login, password string) (*model.User, error)
}

// DocsServiceInterface defines the contract for document service
type DocsServiceInterface interface {
	Create(doc *model.Document) error
	List(owner string, limit int) ([]model.Document, error)
	GetByID(id string) (*model.Document, error)
	Delete(id string) error
}

// SessionServiceInterface defines the contract for session service
type SessionServiceInterface interface {
	Create(userID, login string) string
	Get(token string) (*model.Session, bool)
	Validate(token string) (model.Session, bool)
	Delete(token string) bool
}