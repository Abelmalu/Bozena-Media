package ierrors

import (
	"errors"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAppError_HTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected int
	}{
		{"Validation", TypeValidation, http.StatusBadRequest},
		{"NotFound", TypeNotFound, http.StatusNotFound},
		{"Conflict", TypeConflict, http.StatusConflict},
		{"Unauthorized", TypeUnauthorized, http.StatusUnauthorized},
		{"Forbidden", TypeForbidden, http.StatusForbidden},
		{"Timeout", TypeTimeout, http.StatusGatewayTimeout},
		{"Internal", TypeInternal, http.StatusInternalServerError},
		{"Unknown", "UNKNOWN", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &AppError{Type: tt.errType}
			if got := e.HTTPStatus(); got != tt.expected {
				t.Errorf("AppError.HTTPStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFromGRPC(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorType
	}{
		{"Nil", nil, ""},
		{"InvalidArgument", status.Error(codes.InvalidArgument, "invalid"), TypeValidation},
		{"NotFound", status.Error(codes.NotFound, "not found"), TypeNotFound},
		{"AlreadyExists", status.Error(codes.AlreadyExists, "exists"), TypeConflict},
		{"Unauthenticated", status.Error(codes.Unauthenticated, "unauth"), TypeUnauthorized},
		{"PermissionDenied", status.Error(codes.PermissionDenied, "denied"), TypeForbidden},
		{"DeadlineExceeded", status.Error(codes.DeadlineExceeded, "timeout"), TypeTimeout},
		{"Internal", status.Error(codes.Internal, "internal"), TypeInternal},
		{"NonGRPC", errors.New("normal error"), TypeInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromGRPC(tt.err)
			if tt.err == nil {
				if got != nil {
					t.Errorf("FromGRPC() = %v, want nil", got)
				}
				return
			}
			if got.Type != tt.expected {
				t.Errorf("FromGRPC() type = %v, want %v", got.Type, tt.expected)
			}
		})
	}
}
