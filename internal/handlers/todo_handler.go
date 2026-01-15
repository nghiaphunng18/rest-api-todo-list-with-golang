package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todo-api/internal/services"
	"todo-api/internal/utils"

	"github.com/go-chi/chi"
)

type TodoHandler struct {
    todoService *services.TodoService
}

type CreateTodoRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
}

type UpdateTodoRequest struct {
    Title       *string `json:"title"`
    Description *string `json:"description"`
    Completed   *bool   `json:"completed"`
}

func NewTodoHandler(todoService *services.TodoService) *TodoHandler {
    return &TodoHandler{todoService: todoService}
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
        utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated") 
        return
    }

    // get pagination parameters
    const defaultLimit = 10
    const defaultPage = 1
    
    limit := defaultLimit
    if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }

    page := defaultPage
    if pageStr := r.URL.Query().Get("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }
    
    // call service
    paginatedTodos, err := h.todoService.GetTodos(userID, limit, page)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todos")
        return
    }
    
    response := map[string]interface{}{
        "data": paginatedTodos.Todos,
        "meta": map[string]interface{}{
            "total_items": paginatedTodos.TotalCount,
            "page":        page,
            "limit":       limit,
            "total_pages": (paginatedTodos.TotalCount + limit - 1) / limit, 
        },
    }

    utils.DTOResponse(w, http.StatusOK, response)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    // get user id from context
    userID, err := getUserIDFromContext(r)
    if err != nil {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    var req CreateTodoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // validate
    if req.Title == "" {
        utils.ErrorResponse(w, http.StatusBadRequest, "Title is required")
        return
    }

    todo, err := h.todoService.CreateTodo(userID, req.Title, req.Description)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create todo")
        return
    }

    utils.DTOResponse(w, http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
    // get user id from context
    userID, err := getUserIDFromContext(r)
    if err != nil {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    // get todo id from URL
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid Todo ID")
        return
    }

    todo, err := h.todoService.GetTodoByID(uint(id), userID)
    if err != nil {
        if errors.Is(err, services.ErrTodoNotFound) {
            utils.ErrorResponse(w, http.StatusNotFound, "Todo not found")
        } else {
            utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todo")
        }
        return
    }

    utils.DTOResponse(w, http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
    // get user id from context
    userID, err := getUserIDFromContext(r)
    if err != nil {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    // get id from URL
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid Todo ID")
        return
    }

    var req UpdateTodoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // validate
    if req.Title != nil && *req.Title == "" {
        utils.ErrorResponse(w, http.StatusBadRequest, "Title cannot be empty")
        return
    }

    updatedTodo, err := h.todoService.UpdateTodo(userID, uint(id), req.Title, req.Description, req.Completed)
    if err != nil {
        if errors.Is(err, services.ErrTodoNotFound) {
            utils.ErrorResponse(w, http.StatusNotFound, "Todo not found")
        } else {
            utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update todo")
        }
        return
    }

    utils.DTOResponse(w, http.StatusOK, updatedTodo)
}
