package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func TestCustomError(t *testing.T) {
	// Create a new CustomError
	err := New("TestError", errors.New("caused by error"), "detail message")

	// Check the Error method
	expectedErrorString := fmt.Sprintf("Name: %s, CausedBy: %v\nfile: %s , line: %d",
		err.Name,
		err.CausedBy,
		err.file,
		err.line,
	)
	if err.Error() != expectedErrorString {
		t.Errorf("Expected error string to be '%s', got '%s'", expectedErrorString, err.Error())
	}

	// Check the Unwrap method
	if !errors.Is(err.Unwrap(), err.CausedBy) {
		t.Errorf("Expected unwrapped error to be '%v', got '%v'", err.CausedBy, err.Unwrap())
	}
}
