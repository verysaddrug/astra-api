package mocks

import (
	"astra-api/internal/model"
	"database/sql"
)

type UserRepositoryMock struct {
	CreateFunc     func(user *model.User) error
	CreateTxFunc   func(tx *sql.Tx, user *model.User) error
	GetByLoginFunc func(login string) (*model.User, error)
	GetByIDFunc    func(id string) (*model.User, error)
}

func (m *UserRepositoryMock) Create(user *model.User) error { return m.CreateFunc(user) }
func (m *UserRepositoryMock) CreateTx(tx *sql.Tx, user *model.User) error {
	return m.CreateTxFunc(tx, user)
}
func (m *UserRepositoryMock) GetByLogin(login string) (*model.User, error) {
	return m.GetByLoginFunc(login)
}
func (m *UserRepositoryMock) GetByID(id string) (*model.User, error) { return m.GetByIDFunc(id) }

type DocumentRepositoryMock struct {
	CreateFunc   func(doc *model.Document) error
	CreateTxFunc func(tx *sql.Tx, doc *model.Document) error
	ListFunc     func(owner string, limit int) ([]model.Document, error)
	GetByIDFunc  func(id string) (*model.Document, error)
	DeleteFunc   func(id string) error
	DeleteTxFunc func(tx *sql.Tx, id string) error
}

func (m *DocumentRepositoryMock) Create(doc *model.Document) error { return m.CreateFunc(doc) }
func (m *DocumentRepositoryMock) CreateTx(tx *sql.Tx, doc *model.Document) error {
	return m.CreateTxFunc(tx, doc)
}
func (m *DocumentRepositoryMock) List(owner string, limit int) ([]model.Document, error) {
	return m.ListFunc(owner, limit)
}
func (m *DocumentRepositoryMock) GetByID(id string) (*model.Document, error) {
	return m.GetByIDFunc(id)
}
func (m *DocumentRepositoryMock) Delete(id string) error               { return m.DeleteFunc(id) }
func (m *DocumentRepositoryMock) DeleteTx(tx *sql.Tx, id string) error { return m.DeleteTxFunc(tx, id) }
