package model

import "time"

// Session represents a user session
type Session struct {
	Token   string    `json:"token"`
	UserID  string    `json:"user_id"`
	Login   string    `json:"login"`
	Created time.Time `json:"created"`
}
