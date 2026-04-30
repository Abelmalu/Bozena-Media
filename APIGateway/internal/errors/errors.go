package errors

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
)

type AppError struct {
	Type      ErrorType              `json:"type"`
	Code      string                 `json:"code,omitempty"`
	Message   string                 `json:"message"`
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

func NewAppError(errType ErrorType, message string, cause error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}

func NewValidationError(message string, details map[string]interface{}, cause error) *AppError {
	return &AppError{
		Type:    TypeValidation,
		Message: message,
		Details: details,
		Cause:   cause,
		Code:"VALIDATION_ERROR",
	}
}

func NewNotFoundError(message string, cause error) *AppError {
	return &AppError{
		Type:    TypeNotFound,
		Message: message,
		Cause:   cause,
	}
}

func NewConflictError(message string, cause error) *AppError {
	return &AppError{
		Type:    TypeConflict,
		Message: message,
		Cause:   cause,
	}
}

func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Type:    TypeInternal,
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

func FromGRPC(err error) *AppError {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return NewInternalError("An unexpected error occurred", err)
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return NewValidationError(st.Message(), nil, err)
	case codes.NotFound:
		return NewNotFoundError(st.Message(), err)
	case codes.AlreadyExists:
		return NewConflictError(st.Message(), err)
	case codes.Unauthenticated:
		return NewAppError(TypeUnauthorized, st.Message(), err)
	case codes.PermissionDenied:
		return NewAppError(TypeForbidden, st.Message(), err)
	case codes.DeadlineExceeded:
		return NewAppError(TypeTimeout, "Request timed out", err)
	default:
		return NewInternalError(st.Message(), err)
	}
}
