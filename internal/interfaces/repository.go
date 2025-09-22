package interfaces

import (
	"astra-api/internal/model"
	"database/sql"
)

// UserRepositoryInterface defines the contract for user repository
type UserRepositoryInterface interface {
	Create(user *model.User) error
	CreateTx(tx *sql.Tx, user *model.User) error
	GetByLogin(login string) (*model.User, error)
	GetByID(id string) (*model.User, error)
}

// DocumentRepositoryInterface defines the contract for document repository
type DocumentRepositoryInterface interface {
	Create(doc *model.Document) error
	CreateTx(tx *sql.Tx, doc *model.Document) error
	List(owner string, limit int) ([]model.Document, error)
	GetByID(id string) (*model.Document, error)
	Delete(id string) error
	DeleteTx(tx *sql.Tx, id string) error
}