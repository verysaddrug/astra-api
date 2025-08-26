package repository

import (
	"astra-api/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	_, err := r.db.Exec(`INSERT INTO users (id, login, password, created_at) VALUES ($1, $2, $3, $4)`, user.ID, user.Login, user.Password, user.CreatedAt)
	return err
}

func (r *UserRepository) GetByLogin(login string) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, `SELECT * FROM users WHERE login = $1`, login)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
