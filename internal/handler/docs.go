package handler

import (
	"astra-api/internal/cache"
	"astra-api/internal/model"
	"astra-api/internal/repository"
	"astra-api/internal/service"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DocsHandler struct {
	docsService    service.DocsServiceInterface
	cache          *cache.Cache
	sessionService service.SessionServiceInterface
	userRepo       repository.UserRepositoryInterface
}

func NewDocsHandler(docsService service.DocsServiceInterface, cache *cache.Cache, sessionService service.SessionServiceInterface, userRepo repository.UserRepositoryInterface) *DocsHandler {
	return &DocsHandler{docsService: docsService, cache: cache, sessionService: sessionService, userRepo: userRepo}
}

// @Summary Загрузка документа
// @Tags docs
// @Accept multipart/form-data
// @Produce json
// @Param token query string true "Токен"
// @Param meta formData string true "Метаданные"
// @Param file formData file false "Файл"
// @Success 200 {object} model.APIResponse
// @Router /api/docs [post]
func (h *DocsHandler) Upload(w http.ResponseWriter, r *http.Request) {
	token := GetToken(r)
	sess, ok := h.sessionService.Validate(token)
	if !ok {
		WriteError(w, 401, "invalid token")
		return
	}
	if r.Method != http.MethodPost {
		WriteError(w, 405, "method not allowed")
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		WriteError(w, 400, "invalid multipart form")
		return
	}
	metaStr := r.FormValue("meta")
	var meta struct {
		Name   string   `json:"name"`
		File   bool     `json:"file"`
		Public bool     `json:"public"`
		Mime   string   `json:"mime"`
		Grants []string `json:"grants"`
	}
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		WriteError(w, 400, "invalid meta json")
		return
	}
	var jsonData []byte
	var jsonObj interface{}
	if jsonStr := r.FormValue("json"); jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
			WriteError(w, 400, "invalid json field")
			return
		}
		jsonData = []byte(jsonStr)
	}
	if meta.File {
		file, header, err := r.FormFile("file")
		if err != nil {
			WriteError(w, 400, "file not found in form")
			return
		}
		defer file.Close()
		fileName := filepath.Join("uploads", header.Filename)
		out, err := os.Create(fileName)
		if err != nil {
			WriteError(w, 500, "cannot save file")
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			WriteError(w, 500, "cannot write file")
			return
		}
	}
	doc := &model.Document{
		Name:     meta.Name,
		Mime:     meta.Mime,
		File:     meta.File,
		Public:   meta.Public,
		Owner:    sess.UserID,
		Grants:   meta.Grants,
		JsonData: jsonData,
	}
	if err := h.docsService.Create(doc); err != nil {
		WriteError(w, 500, err.Error())
		return
	}
	h.cache.InvalidateAll()
	resp := map[string]interface{}{
		"id":   doc.ID,
		"file": meta.Name,
	}
	if jsonObj != nil {
		resp["json"] = jsonObj
	} else {
		resp["json"] = nil
	}
	WriteResponse(w, &model.APIResponse{Data: resp})
}

// @Summary Список документов
// @Tags docs
// @Produce json
// @Param token query string true "Токен"
// @Param login query string false "Логин"
// @Param limit query int false "Лимит"
// @Success 200 {object} model.APIResponse
// @Router /api/docs [get]
func (h *DocsHandler) List(w http.ResponseWriter, r *http.Request) {
	token := GetToken(r)
	sess, ok := h.sessionService.Validate(token)
	if !ok {
		WriteError(w, 401, "invalid token")
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		WriteError(w, 405, "method not allowed")
		return
	}
	ownerID := sess.UserID
	if login := r.URL.Query().Get("login"); login != "" {
		user, err := h.userRepo.GetByLogin(login)
		if err != nil {
			WriteError(w, 400, "unknown login")
			return
		}
		ownerID = user.ID
	}
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	cacheKey := "list:" + ownerID + ":" + strconv.Itoa(limit)
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		if cached, ok := h.cache.Get(cacheKey); ok {
			WriteResponse(w, &model.APIResponse{Data: map[string]interface{}{"docs": cached}})
			return
		}
	}
	docs, err := h.docsService.List(ownerID, limit)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
	}
	h.cache.Set(cacheKey, docs)
	WriteResponse(w, &model.APIResponse{Data: map[string]interface{}{"docs": docs}})
}

// @Summary Документ по id
// @Tags docs
// @Produce json
// @Param token query string true "Токен"
// @Param id path string true "ID"
// @Success 200 {object} model.APIResponse
// @Router /api/docs/{id} [get]
func (h *DocsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	token := GetToken(r)
	if _, ok := h.sessionService.Validate(token); !ok {
		WriteError(w, 401, "invalid token")
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		WriteError(w, 405, "method not allowed")
		return
	}
	id := getIDFromURL(r.URL.Path)
	if id == "" {
		WriteError(w, 400, "missing document id")
		return
	}

	// Сначала пробуем из кэша
	cacheKey := "doc:" + id
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		if cached, ok := h.cache.Get(cacheKey); ok {
			doc := cached.(*model.Document)
			if doc.File {
				filePath := "uploads/" + doc.Name
				f, err := os.Open(filePath)
				if err != nil {
					WriteError(w, 404, "file not found")
					return
				}
				defer f.Close()
				w.Header().Set("Content-Type", doc.Mime)
				w.Header().Set("Content-Disposition", "attachment; filename="+doc.Name)
				if r.Method == http.MethodGet {
					io.Copy(w, f)
				}
				return
			}
			var obj interface{}
			if len(doc.JsonData) > 0 {
				_ = json.Unmarshal(doc.JsonData, &obj)
			}
			WriteResponse(w, &model.APIResponse{Data: obj})
			return
		}
	}

	// Если не в кэше - ищем в БД
	doc, err := h.docsService.GetByID(id)
	if err != nil {
		WriteError(w, 404, "document not found")
		return
	}

	// Сохраняем в кэш
	h.cache.Set(cacheKey, doc)

	if doc.File {
		filePath := "uploads/" + doc.Name
		f, err := os.Open(filePath)
		if err != nil {
			WriteError(w, 404, "file not found")
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", doc.Mime)
		w.Header().Set("Content-Disposition", "attachment; filename="+doc.Name)
		if r.Method == http.MethodGet {
			io.Copy(w, f)
		}
		return
	}

	var obj interface{}
	if len(doc.JsonData) > 0 {
		_ = json.Unmarshal(doc.JsonData, &obj)
	}
	WriteResponse(w, &model.APIResponse{Data: obj})
}

// @Summary Удалить документ
// @Tags docs
// @Produce json
// @Param token query string true "Токен"
// @Param id path string true "ID"
// @Success 200 {object} model.APIResponse
// @Router /api/docs/{id} [delete]
func (h *DocsHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	token := GetToken(r)
	if _, ok := h.sessionService.Validate(token); !ok {
		WriteError(w, 401, "invalid token")
		return
	}
	if r.Method != http.MethodDelete {
		WriteError(w, 405, "method not allowed")
		return
	}
	id := getIDFromURL(r.URL.Path)
	if id == "" {
		WriteError(w, 400, "missing document id")
		return
	}
	if err := h.docsService.Delete(id); err != nil {
		WriteError(w, 404, "document not found")
		return
	}
	h.cache.Invalidate("doc:" + id)
	h.cache.InvalidateAll()
	WriteResponse(w, &model.APIResponse{Response: map[string]bool{id: true}})
}

func getIDFromURL(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return ""
	}
	return parts[3]
}
