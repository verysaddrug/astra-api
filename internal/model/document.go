package model

import (
	"time"

	"github.com/lib/pq"
)

type Document struct {
	ID        string         `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Mime      string         `db:"mime" json:"mime"`
	File      bool           `db:"file" json:"file"`
	Public    bool           `db:"public" json:"public"`
	Owner     string         `db:"owner" json:"owner"`
	CreatedAt time.Time      `db:"created_at" json:"created"`
	Grants    pq.StringArray `db:"grants" json:"grants"`
	JsonData  []byte         `db:"json_data" json:"json,omitempty"`
}
