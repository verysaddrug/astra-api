package mocks

import (
	"astra-api/internal/model"
)

type AuthServiceMock struct {
	RegisterFunc     func(login, password, adminToken string) (*model.User, error)
	AuthenticateFunc func(login, password string) (*model.User, error)
}

func (m *AuthServiceMock) Register(login, password, adminToken string) (*model.User, error) {
	return m.RegisterFunc(login, password, adminToken)
}

func (m *AuthServiceMock) Authenticate(login, password string) (*model.User, error) {
	return m.AuthenticateFunc(login, password)
}

type DocsServiceMock struct {
	CreateFunc  func(doc *model.Document) error
	ListFunc    func(owner string, limit int) ([]model.Document, error)
	GetByIDFunc func(id string) (*model.Document, error)
	DeleteFunc  func(id string) error
}

func (m *DocsServiceMock) Create(doc *model.Document) error { return m.CreateFunc(doc) }
func (m *DocsServiceMock) List(owner string, limit int) ([]model.Document, error) {
	return m.ListFunc(owner, limit)
}
func (m *DocsServiceMock) GetByID(id string) (*model.Document, error) { return m.GetByIDFunc(id) }
func (m *DocsServiceMock) Delete(id string) error                     { return m.DeleteFunc(id) }

type SessionServiceMock struct {
	CreateFunc   func(userID, login string) string
	GetFunc      func(token string) (*model.Session, bool)
	ValidateFunc func(token string) (model.Session, bool)
	DeleteFunc   func(token string) bool
}

func (m *SessionServiceMock) Create(userID, login string) string      { return m.CreateFunc(userID, login) }
func (m *SessionServiceMock) Get(token string) (*model.Session, bool) { return m.GetFunc(token) }
func (m *SessionServiceMock) Validate(token string) (model.Session, bool) {
	return m.ValidateFunc(token)
}
func (m *SessionServiceMock) Delete(token string) bool { return m.DeleteFunc(token) }
