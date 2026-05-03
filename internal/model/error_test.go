package model

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewApiError(t *testing.T) {
	inputError := errors.New("test error")
	inputCode := http.StatusNotFound
	apiError := NewApiError(inputError, inputCode)
	
	if apiError.Code != inputCode {
		t.Errorf("expected code %d, got %d", inputCode, apiError.Code)
	}

	if apiError.Message != inputError.Error() {
		t.Errorf("expected message %s, got %s", inputError.Error(), apiError.Message)
	}

	if apiError.Err.Error() != inputError.Error() {
		t.Errorf("expected underlying error %v, got %v", inputError, apiError.Err)
	}
}

func TestApiError_Error(t *testing.T) {
	msg := "database failure"
	apiError := &ApiError{
		Err: errors.New(msg),
	}

	expected := msg
	if apiError.Error() != expected {
		t.Errorf("expected error string %s, got %s", expected, apiError.Error())
	}
}
