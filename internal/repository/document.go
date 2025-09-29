package repository

import (
	"astra-api/internal/model"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DocumentRepository struct {
	db *sqlx.DB
}

func NewDocumentRepository(db *sqlx.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(doc *model.Document) error {
	var jsonArg interface{}
	if len(doc.JsonData) > 0 {
		jsonArg = string(doc.JsonData)
	} else {
		jsonArg = nil
	}
	_, err := r.db.Exec(`INSERT INTO documents (id, name, mime, file, public, owner, created_at, grants, json_data) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb)`, doc.ID, doc.Name, doc.Mime, doc.File, doc.Public, doc.Owner, doc.CreatedAt, pq.Array(doc.Grants), jsonArg)
	return err
}

func (r *DocumentRepository) CreateTx(tx *sql.Tx, doc *model.Document) error {
	var jsonArg interface{}
	if len(doc.JsonData) > 0 {
		jsonArg = string(doc.JsonData)
	} else {
		jsonArg = nil
	}
	_, err := tx.Exec(`INSERT INTO documents (id, name, mime, file, public, owner, created_at, grants, json_data) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb)`, doc.ID, doc.Name, doc.Mime, doc.File, doc.Public, doc.Owner, doc.CreatedAt, pq.Array(doc.Grants), jsonArg)
	return err
}

func (r *DocumentRepository) List(owner string, limit int) ([]model.Document, error) {
	docs := []model.Document{}
	err := r.db.Select(&docs, `SELECT id, name, mime, file, public, owner, created_at, grants::text[] as grants, json_data FROM documents WHERE owner = $1 ORDER BY name, created_at DESC LIMIT $2`, owner, limit)
	return docs, err
}

func (r *DocumentRepository) GetByID(id string) (*model.Document, error) {
	var doc model.Document
	err := r.db.Get(&doc, `SELECT id, name, mime, file, public, owner, created_at, grants::text[] as grants, json_data FROM documents WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *DocumentRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM documents WHERE id = $1`, id)
	return err
}

func (r *DocumentRepository) DeleteTx(tx *sql.Tx, id string) error {
	_, err := tx.Exec(`DELETE FROM documents WHERE id = $1`, id)
	return err
}
