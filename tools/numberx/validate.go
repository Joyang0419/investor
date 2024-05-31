package numberx

import (
	"golang.org/x/exp/constraints"
)

type TypeNumber interface {
	constraints.Float | constraints.Integer | constraints.Signed | constraints.Unsigned
}

func IsNotZero[typeNumber TypeNumber](number typeNumber) bool {
	return number != typeNumber(0)
}

func IsGTE[typeNumber TypeNumber](number typeNumber, compareNumber typeNumber) bool {
	return number >= compareNumber
}

func IsLTE[typeNumber TypeNumber](number typeNumber, compareNumber typeNumber) bool {
	return number <= compareNumber
}

func IsLT[typeNumber TypeNumber](number typeNumber, compareNumber typeNumber) bool {
	return number < compareNumber
}

func IsGT[typeNumber TypeNumber](number typeNumber, compareNumber typeNumber) bool {
	return number > compareNumber
}
