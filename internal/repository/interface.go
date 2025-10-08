package repository

import (
	"astra-api/internal/model"
	"database/sql"
)

// UserRepositoryInterface описывает контракт репозитория пользователей
type UserRepositoryInterface interface {
	Create(user *model.User) error
	CreateTx(tx *sql.Tx, user *model.User) error
	GetByLogin(login string) (*model.User, error)
	GetByID(id string) (*model.User, error)
}

// DocumentRepositoryInterface описывает контракт репозитория документов
type DocumentRepositoryInterface interface {
	Create(doc *model.Document) error
	CreateTx(tx *sql.Tx, doc *model.Document) error
	List(owner string, limit int) ([]model.Document, error)
	GetByID(id string) (*model.Document, error)
	Delete(id string) error
	DeleteTx(tx *sql.Tx, id string) error
}
