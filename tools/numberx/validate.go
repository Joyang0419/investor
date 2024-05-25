package numberx

import (
	"golang.org/x/exp/constraints"
)

func IsNotZero[numberType constraints.Float | constraints.Integer | constraints.Signed | constraints.Unsigned](number numberType) bool {
	return number != numberType(0)
}
