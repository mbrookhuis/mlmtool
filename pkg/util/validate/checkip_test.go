package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unit test for CheckIP function
func TestCheckIP(t *testing.T) {
	tests := []struct {
		ip        string
		expectErr bool
	}{
		{"192.168.1.1", false},
		{"255.255.255.255", false},
		{"0.0.0.0", false},
		{"999.999.999.999", true},
		{"random-string", true},
	}

	for _, tt := range tests {
		err := CheckIP(tt.ip)
		if tt.expectErr {
			assert.Error(t, err, "Expected no error for valid IP: %v", err)
		} else {
			assert.NoError(t, err, "Expected error for invalid IP: %v", err)
		}
	}
}
