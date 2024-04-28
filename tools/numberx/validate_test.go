package numberx_test

import (
	"testing"

	"tools/numberx"
)

func TestCheckNotZeroWithPositiveNumber(t *testing.T) {
	if !numberx.CheckNotZero[int](5) {
		t.Errorf("CheckNotZero failed with positive number")
	}
}

func TestCheckNotZeroWithNegativeNumber(t *testing.T) {
	if !numberx.CheckNotZero[int](-5) {
		t.Errorf("CheckNotZero failed with negative number")
	}
}

func TestCheckNotZeroWithZero(t *testing.T) {
	if numberx.CheckNotZero[int](0) {
		t.Errorf("CheckNotZero failed with zero")
	}
}

func TestCheckNotZeroWithPositiveFloat(t *testing.T) {
	if !numberx.CheckNotZero[float64](5.5) {
		t.Errorf("CheckNotZero failed with positive float")
	}
}

func TestCheckNotZeroWithNegativeFloat(t *testing.T) {
	if !numberx.CheckNotZero[float64](-5.5) {
		t.Errorf("CheckNotZero failed with negative float")
	}
}

func TestCheckNotZeroWithZeroFloat(t *testing.T) {
	if numberx.CheckNotZero[float64](0.0) {
		t.Errorf("CheckNotZero failed with zero float")
	}
}
