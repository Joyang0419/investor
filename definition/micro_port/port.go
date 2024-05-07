package micro_port

import (
	"fmt"
)

const (
	MicroStockPricePort = 50051
	MicroAuthPort       = 50052
)

func GetGrpcAddress(domain string, port int) string {
	return fmt.Sprintf("%s:%d", domain, port)
}
