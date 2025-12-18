package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"todo-api/internal/services"
	"todo-api/internal/utils"
)

type TodoHandler struct {
    todoService *services.TodoService
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