package ierrors

import (
	"errors"
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
	TypeTooManyRequests ErrorType = "TOO_MANY_REQUESTS"
)

var(

	ErrUserIDNotFoundInContext error = errors.New("user id not found in context")
	ErrTypeAssertionFailed error = errors.New("type assertion failed")
	ErrJTINotFoundInContext error = errors.New("JTI not found in context")
	ErrExpTimeNotFoundInContext error = errors.New("expiration time not found in context")
	ErrRequestIDNotFoundInContext error = errors.New("request id not found in context")
	ErrMetaDataNotFound = errors.New("meta data not found")
	
)


type AppError struct {
	Type    ErrorType              `json:"type"`
	Code    string                 `json:"code,omitempty"`
	Message ErrorMessage           `json:"message"`
	Cause   error                  `json:"-"`
	Details map[string]interface{} `json:"details,omitempty"`
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

func NewTooManyRequestsError(message ErrorMessage, details map[string]interface{}, cause error) *AppError {
	return &AppError{
		Type:    TypeTooManyRequests,
		Message: message,
		Details: details,
		Cause:   cause,
		Code:    "TOO_MANY_REQUESTS",
	}
}

func NewInternalError(message ErrorMessage, cause error) *AppError {
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
	case TypeTooManyRequests:
		return http.StatusTooManyRequests
		
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
		errStr := st.Message()
		return NewValidationError(ErrorMessage(errStr), nil, err)
	case codes.NotFound:
		errStr := st.Message()
		return NewNotFoundError(ErrorMessage(errStr), err)
	case codes.AlreadyExists:
		errStr := st.Message()

		return NewConflictError(ErrorMessage(errStr), err)
	case codes.Unauthenticated:
		errStr := st.Message()
		return NewAppError(TypeUnauthorized, ErrorMessage(errStr), err)
	case codes.PermissionDenied:
		errStr := st.Message()
		return NewAppError(TypeForbidden, ErrorMessage(errStr), err)
	case codes.DeadlineExceeded:
		errStr := st.Message()
		return NewAppError(TypeTimeout, ErrorMessage(errStr), err)
	default:
		return NewInternalError(MSGSomethingWentWrong, err)
	}
}
