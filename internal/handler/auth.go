package handler

import (
	"astra-api/internal/model"
	"astra-api/internal/service"
	"encoding/json"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService    *service.AuthService
	sessionService *service.SessionService
}

func NewAuthHandler(authService *service.AuthService, sessionService *service.SessionService) *AuthHandler {
	return &AuthHandler{authService: authService, sessionService: sessionService}
}

// @Summary Регистрация
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.RegisterRequest true "Данные"
// @Success 200 {object} model.APIResponse
// @Router /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, 405, "method not allowed")
		return
	}
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, 400, "invalid request body")
		return
	}
	user, err := h.authService.Register(req.Login, req.Pswd, req.Token)
	if err != nil {
		WriteError(w, 400, err.Error())
		return
	}
	WriteResponse(w, &model.APIResponse{Response: map[string]string{"login": user.Login}})
}

// @Summary Логин
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.AuthRequest true "Данные"
// @Success 200 {object} model.APIResponse
// @Router /api/auth [post]
func (h *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, 405, "method not allowed")
		return
	}
	var req model.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, 400, "invalid request body")
		return
	}
	user, err := h.authService.Authenticate(req.Login, req.Pswd)
	if err != nil {
		WriteError(w, 401, err.Error())
		return
	}
	token := h.sessionService.Create(user.ID, user.Login)
	WriteResponse(w, &model.APIResponse{Response: map[string]string{"token": token}})
}

// @Summary Логаут
// @Tags auth
// @Produce json
// @Param token path string true "Токен"
// @Success 200 {object} model.APIResponse
// @Router /api/auth/{token} [delete]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteError(w, 405, "method not allowed")
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		WriteError(w, 400, "missing token")
		return
	}
	token := parts[3]
	ok := h.sessionService.Delete(token)
	if !ok {
		WriteError(w, 400, "invalid token")
		return
	}
	WriteResponse(w, &model.APIResponse{Response: map[string]bool{token: true}})
}

func WriteError(w http.ResponseWriter, code int, text string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&model.APIResponse{Error: &model.APIError{Code: code, Text: text}})
}

func WriteResponse(w http.ResponseWriter, resp *model.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}
