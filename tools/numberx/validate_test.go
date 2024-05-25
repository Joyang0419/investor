package numberx_test

import (
	"testing"

	"tools/numberx"
)

func TestCheckNotZeroWithPositiveNumber(t *testing.T) {
	if !numberx.IsNotZero[int](5) {
		t.Errorf("IsNotZero failed with positive number")
	}
}

func TestCheckNotZeroWithNegativeNumber(t *testing.T) {
	if !numberx.IsNotZero[int](-5) {
		t.Errorf("IsNotZero failed with negative number")
	}
}

func TestCheckNotZeroWithZero(t *testing.T) {
	if numberx.IsNotZero[int](0) {
		t.Errorf("IsNotZero failed with zero")
	}
}

func TestCheckNotZeroWithPositiveFloat(t *testing.T) {
	if !numberx.IsNotZero[float64](5.5) {
		t.Errorf("IsNotZero failed with positive float")
	}
}

func TestCheckNotZeroWithNegativeFloat(t *testing.T) {
	if !numberx.IsNotZero[float64](-5.5) {
		t.Errorf("IsNotZero failed with negative float")
	}
}

func TestCheckNotZeroWithZeroFloat(t *testing.T) {
	if numberx.IsNotZero[float64](0.0) {
		t.Errorf("IsNotZero failed with zero float")
	}
}
