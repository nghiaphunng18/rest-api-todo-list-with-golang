package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"todo-api/internal/services"
	"todo-api/internal/utils"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
        return
    }

    if req.Username == "" || req.Password == "" {
        http.Error(w, `{"error": "missing username or password"}`, http.StatusBadRequest)
        return
    }

	if err := h.authService.Register(req.Username, req.Password); err != nil {
		if errors.Is(err, services.ErrUsernameExists) {
			utils.ErrorResponse(w, http.StatusConflict, "Username already exists")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to register user")
		return
	}
	

	utils.DTOResponse(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
    }

    token, err := h.authService.Login(req.Username, req.Password)
    if err != nil {
        utils.ErrorResponse(w, http.StatusUnauthorized, "invalid credentials")
		return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": token})
}