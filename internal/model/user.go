package model

import "time"

// User пользователь системы
// @Description Модель пользователя
// @Description Пароль хранится в хеше
// @Description Логин уникален
// @Description В ответах пароль не возвращается
// @Description created_at — дата создания
// @Description id — UUID
// @Description login — строка
// @Description password — строка (hash)
// @Description created_at — строка (timestamp)
type User struct {
	ID        string    `db:"id" json:"id" example:"b1a7c8e2-1c2d-4e5f-8a7b-2c3d4e5f6a7b"`
	Login     string    `db:"login" json:"login" example:"TestUser01"`
	Password  string    `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at" example:"2024-08-26T10:30:56Z"`
}

type RegisterRequest struct {
	Token string `json:"token" example:"supersecrettoken"`
	Login string `json:"login" example:"TestUser01"`
	Pswd  string `json:"pswd" example:"Qwerty123!"`
}

type AuthRequest struct {
	Login string `json:"login" example:"TestUser01"`
	Pswd  string `json:"pswd" example:"Qwerty123!"`
}
