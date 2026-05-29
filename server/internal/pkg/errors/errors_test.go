package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(1001, 400, "bad request")
	if err.Code != 1001 {
		t.Errorf("expected code 1001, got %d", err.Code)
	}
	if err.Status != 400 {
		t.Errorf("expected status 400, got %d", err.Status)
	}
	if err.Message != "bad request" {
		t.Errorf("expected message 'bad request', got %s", err.Message)
	}
	if err.Internal() != nil {
		t.Error("expected nil internal error")
	}
}

func TestWrap(t *testing.T) {
	inner := errors.New("db connection failed")
	err := Wrap(inner, 1006, 500, "internal error")

	if !errors.Is(err, inner) {
		t.Error("expected errors.Is to match inner error")
	}
	if err.Internal() != inner {
		t.Error("expected internal error to be inner")
	}
}

func TestWrapInternal(t *testing.T) {
	inner := errors.New("panic")
	err := WrapInternal(inner)

	if err.Code != ErrInternal.Code {
		t.Errorf("expected code %d, got %d", ErrInternal.Code, err.Code)
	}
	if err.Status != 500 {
		t.Errorf("expected status 500, got %d", err.Status)
	}
}

func TestWithMessage(t *testing.T) {
	err := New(1001, 400, "original")
	newErr := WithMessage(err, "modified")

	if newErr.Message != "modified" {
		t.Errorf("expected message 'modified', got %s", newErr.Message)
	}
	if newErr.Code != err.Code {
		t.Error("expected code unchanged")
	}
}

func TestWithStatus(t *testing.T) {
	err := New(1001, 400, "test")
	newErr := WithStatus(err, 422)

	if newErr.Status != 422 {
		t.Errorf("expected status 422, got %d", newErr.Status)
	}
}

func TestIsAppError(t *testing.T) {
	appErr := New(1004, 404, "not found")
	stdErr := errors.New("standard error")

	if _, ok := IsAppError(stdErr); ok {
		t.Error("expected IsAppError to return false for standard error")
	}

	found, ok := IsAppError(appErr)
	if !ok {
		t.Error("expected IsAppError to return true for AppError")
	}
	if found.Code != 1004 {
		t.Errorf("expected code 1004, got %d", found.Code)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		err    *AppError
		code   int
		status int
	}{
		{ErrBadRequest, 1001, 400},
		{ErrUnauthorized, 1002, 401},
		{ErrForbidden, 1003, 403},
		{ErrNotFound, 1004, 404},
		{ErrConflict, 1005, 409},
		{ErrInternal, 1006, 500},
		{ErrNotImplemented, 1007, 501},
		{ErrTooManyRequests, 1008, 429},
	}

	for _, tt := range tests {
		if tt.err.Code != tt.code {
			t.Errorf("%s: expected code %d, got %d", tt.err.Message, tt.code, tt.err.Code)
		}
		if tt.err.Status != tt.status {
			t.Errorf("%s: expected status %d, got %d", tt.err.Message, tt.status, tt.err.Status)
		}
	}
}
