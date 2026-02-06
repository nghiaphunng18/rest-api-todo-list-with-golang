package domain

import "net/http"

type AppError struct {
    StatusCode int    `json:"-"`
    Message    string `json:"message"`
    ErrorCode  string `json:"error_code"`
}

func (e *AppError) Error() string {
    return e.Message
}

// Define common errors
var (
    ErrInternalServer = &AppError{
        StatusCode: http.StatusInternalServerError,
        Message:    "Internal server error",
        ErrorCode:  "INTERNAL_ERROR",
    }

    ErrInvalidInput = &AppError{
        StatusCode: http.StatusBadRequest,
        Message:    "Invalid input data",
        ErrorCode:  "INVALID_INPUT",
    }

    ErrUnauthorized = &AppError{
        StatusCode: http.StatusUnauthorized,
        Message:    "Unauthorized access",
        ErrorCode:  "UNAUTHORIZED",
    }

    // Auth Errors
    ErrUserExists = &AppError{
        StatusCode: http.StatusConflict,
        Message:    "User already exists",
        ErrorCode:  "USER_ALREADY_EXISTS",
    }

    ErrBadCredentials = &AppError{
        StatusCode: http.StatusUnauthorized,
        Message:    "Invalid username or password",
        ErrorCode:  "INVALID_CREDENTIALS",
    }

    // Todo Errors
    ErrTodoNotFound = &AppError{
        StatusCode: http.StatusNotFound,
        Message:    "Todo not found",
        ErrorCode:  "TODO_NOT_FOUND",
    }
)
