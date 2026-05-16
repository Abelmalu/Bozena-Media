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
	TypeCancelled    ErrorType = "CANCELLED"
	TypeDatabase     ErrorType = "DATABASE"
)

//constant errors 
var (

	ErrMetaDataNotFound = errors.New("meta data not found")
	ErrRequestIDNotFound = errors.New("reques id not found")
)

type AppError struct {
	Type      ErrorType              `json:"type"`
	Code      string                 `json:"code"`
	Message   ErrorMessage           `json:"message"`
	Detail    map[string]interface{} `json:"detail,omitempty"`
	RequestID string                 `json:"request_id"`
	Cause     error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) UnWrap() error {

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
		var MSGInvalidArgument ErrorMessage = ErrorMessage(errStr)
		return NewValidationError(MSGInvalidArgument, nil, err)
	case codes.NotFound:
		errStr := st.Message()
		var MSGNotFound ErrorMessage = ErrorMessage(errStr)
		return NewNotFoundError(MSGNotFound, err)
	case codes.AlreadyExists:
		errStr := st.Message()
		var MSGAreadyExists ErrorMessage = ErrorMessage(errStr)

		return NewConflictError(MSGAreadyExists, err)
	case codes.Unauthenticated:
		errStr := st.Message()
		var MSGUnauthenticated ErrorMessage = ErrorMessage(errStr)
		return NewAppError(TypeUnauthorized, MSGUnauthenticated, err)
	case codes.PermissionDenied:
		errStr := st.Message()
		var MSGPermissionDenied ErrorMessage = ErrorMessage(errStr)
		return NewAppError(TypeForbidden, MSGPermissionDenied, err)
	case codes.DeadlineExceeded:
		errStr := st.Message()
		var MSGDeadlineExceeded ErrorMessage = ErrorMessage(errStr)
		return NewAppError(TypeTimeout, MSGDeadlineExceeded, err)
	default:
		return NewInternalError(MSGSomethingWentWrong, err)
	}
}
