package handlers

import (
	"encoding/json"
	"net/http"
	"todo-api/internal/domain"
	"todo-api/internal/services"
	"todo-api/internal/utils"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService services.IAuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService services.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ParseError(w, domain.ErrInvalidInput)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authService.Register(req.Username, req.Password); err != nil {
		utils.ParseError(w, err)
		return
	}

	utils.DTOResponse(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ParseError(w, domain.ErrInvalidInput)
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.ParseError(w, err)
		return
	}

	utils.DTOResponse(w, http.StatusOK, map[string]string{"token": token})
}