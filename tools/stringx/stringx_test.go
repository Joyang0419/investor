package stringx

import (
	"testing"
)

func TestCheckEmptyReturnsTrueForEmptyString(t *testing.T) {
	result := IsEmptyStr("")
	if result != true {
		t.Errorf("Expected true, got %v", result)
	}
}

func TestCheckEmptyReturnsFalseForNonEmptyString(t *testing.T) {
	result := IsEmptyStr("non-empty")
	if result != false {
		t.Errorf("Expected false, got %v", result)
	}
}

func TestCheckEmptyReturnsTrueForWhitespaceString(t *testing.T) {
	result := IsEmptyStr(" ")
	if result != false {
		t.Errorf("Expected false, got %v", result)
	}
}
