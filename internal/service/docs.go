package service

import (
	"astra-api/internal/interfaces"
	"astra-api/internal/model"
	"time"

	"github.com/google/uuid"
)

type DocsService struct {
	docRepo interfaces.DocumentRepositoryInterface
}

func NewDocsService(docRepo interfaces.DocumentRepositoryInterface) *DocsService {
	return &DocsService{docRepo: docRepo}
}

func (s *DocsService) Create(doc *model.Document) error {
	doc.ID = uuid.New().String()
	doc.CreatedAt = time.Now()
	return s.docRepo.Create(doc)
}

func (s *DocsService) List(owner string, limit int) ([]model.Document, error) {
	return s.docRepo.List(owner, limit)
}

func (s *DocsService) GetByID(id string) (*model.Document, error) {
	return s.docRepo.GetByID(id)
}

func (s *DocsService) Delete(id string) error {
	return s.docRepo.Delete(id)
}
