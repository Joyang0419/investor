package grpcx

import (
	"fmt"
)

func GetGrpcAddress(domain string, port int) string {
	return fmt.Sprintf("%s:%d", domain, port)
}
