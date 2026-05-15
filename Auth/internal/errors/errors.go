package ierrors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	TypeValidation   ErrorType = "VALIDATION"
	TypeNotFound     ErrorType = "NOT_FOUND"
	TypeConflict     ErrorType = "CONFLICT"
	TypeUnauthorized ErrorType = "UNAUTHORIZED"
	TypeForbidden    ErrorType = "FORBIDDEN"
	TypeInternal     ErrorType = "INTERNAL"
	TypeTimeout      ErrorType = "TIMEOUT"
	TypeCancelled    ErrorType = "CANCELLED"
	TypeDatabase     ErrorType = "DATABASE"
)

type AppError struct {
	Type      ErrorType              `json:"type"`
	Code      string                 `json:"code,omitempty"`
	Message   ErrorMessage           `json:"message"`
	Cause     error                  `json:"-"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}
func NewAppError(errType ErrorType, message ErrorMessage, cause error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}

func NewValidationError(message ErrorMessage, details map[string]interface{}, cause error) *AppError {
	return &AppError{
		Type:    TypeValidation,
		Message: message,
		Details: details,
		Cause:   cause,
		Code:    "VALIDATION_ERROR",
	}
}

func NewNotFoundError(message ErrorMessage, cause error) *AppError {
	return &AppError{
		Type:    TypeNotFound,
		Message: message,
		Cause:   cause,
	}
}

func NewConflictError(message ErrorMessage, cause error) *AppError {
	return &AppError{
		Type:    TypeConflict,
		Message: message,
		Cause:   cause,
	}
}

func NewUnauthorizedError(message ErrorMessage, cause error) *AppError {

	return &AppError{
		Type:    TypeUnauthorized,
		Message: message,
		Cause:   cause,
	}
}

func NewInternalError(message ErrorMessage, cause error) *AppError {
	return &AppError{
		Type:    TypeInternal,
		Message: message,
		Cause:   cause,
	}
}

func NewTimeoutError(message ErrorMessage, cause error) *AppError {

	return &AppError{
		Type:    TypeTimeout,
		Message: message,
		Cause:   cause,
	}
}

func NewCancelationError(message ErrorMessage, cause error) *AppError {

	return &AppError{
		Type:    TypeCancelled,
		Message: message,
		Cause:   cause,
	}

}

func NewDatabaseError(message ErrorMessage, cause error) *AppError {

	return &AppError{
		Type:    TypeDatabase,
		Message: message,
		Cause:   cause,
	}
}

func (e *AppError) HTTPStatus() int {
	switch e.Type {
	case TypeValidation:
		return http.StatusBadRequest
	case TypeNotFound:
		return http.StatusNotFound
	case TypeConflict:
		return http.StatusConflict
	case TypeUnauthorized:
		return http.StatusUnauthorized
	case TypeForbidden:
		return http.StatusForbidden
	case TypeTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}
