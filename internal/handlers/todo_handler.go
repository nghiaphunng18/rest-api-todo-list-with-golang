package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo-api/internal/domain"
	"todo-api/internal/services"
	"todo-api/internal/utils"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type TodoHandler struct {
    todoService services.ITodoService
    validate    *validator.Validate
}

type CreateTodoRequest struct {
    Title       string `json:"title" validate:"required,min=3,max=100"`
    Description string `json:"description" validate:"max=500"`
}

type UpdateTodoRequest struct {
    Title       *string `json:"title" validate:"omitempty,min=3,max=100"`
    Description *string `json:"description" validate:"omitempty,max=500"`
    Completed   *bool   `json:"completed"`
}

func NewTodoHandler(todoService services.ITodoService) *TodoHandler {
    return &TodoHandler{
        todoService: todoService,
        validate:    validator.New(),
    }
}

func getUserIDFromContext(r *http.Request) (uint, error) {
	userIDVal := r.Context().Value(utils.UserIDKey)
	
	userID, ok := userIDVal.(uint)
	if !ok {
		return 0, fmt.Errorf("user ID not found") 
	}
	return userID, nil
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
    // exact user id from context
    userID, err := getUserIDFromContext(r)
    if err != nil {
        utils.ParseError(w, domain.ErrUnauthorized)
		return
    }

    // get pagination parameters
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 { limit = 10 }
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 { page = 1 }

	paginated, err := h.todoService.GetTodos(userID, limit, page)
	if err != nil {
		utils.ParseError(w, err)
		return
	}
    
    utils.DTOResponse(w, http.StatusOK, map[string]interface{}{
		"data": paginated.Todos,
		"meta": map[string]interface{}{
			"total_items": paginated.TotalCount,
			"page":        page,
			"limit":       limit,
		},
	})
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    // get user id from context
    userID, err := getUserIDFromContext(r)
    var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ParseError(w, domain.ErrInvalidInput)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.todoService.CreateTodo(userID, req.Title, req.Description)
	if err != nil {
		utils.ParseError(w, err)
		return
	}
	utils.DTOResponse(w, http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
    userID, _ := getUserIDFromContext(r)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	todo, err := h.todoService.GetTodoByID(uint(id), userID)
	if err != nil {
		utils.ParseError(w, err)
		return
	}
	utils.DTOResponse(w, http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
    // get user id from context
    userID, _ := getUserIDFromContext(r)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ParseError(w, domain.ErrInvalidInput)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.todoService.UpdateTodo(userID, uint(id), req.Title, req.Description, req.Completed)
	if err != nil {
		utils.ParseError(w, err)
		return
	}
	utils.DTOResponse(w, http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
    // get user id from context
    userID, _ := getUserIDFromContext(r)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	if err := h.todoService.DeleteTodo(userID, uint(id)); err != nil {
		utils.ParseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
