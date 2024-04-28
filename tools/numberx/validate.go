package numberx

import (
	"golang.org/x/exp/constraints"
)

type TypeNumber interface {
	constraints.Integer | constraints.Float
}

func CheckNotZero[numberType TypeNumber](number numberType) bool {
	return number != numberType(0)
}
