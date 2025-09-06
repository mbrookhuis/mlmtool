package validate

import (
	"fmt"
	"net"
)

// validates, IP is in ip format
// This is a wrapper for net.ParseIP just to provide a custom error message
func CheckIP(IP string) error {
	if net.ParseIP(IP) == nil {
		return fmt.Errorf("ip %s is not correct format", IP)
	}
	return nil
}
