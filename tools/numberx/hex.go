package numberx

import (
	"fmt"
)

func ToCapitalHex[typeNumber TypeNumber](number typeNumber) string {
	return fmt.Sprintf("%X", number)
}
