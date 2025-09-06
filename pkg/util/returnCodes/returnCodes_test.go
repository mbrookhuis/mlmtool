package returncodes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReturnCodes tests the return codes for expected values
func TestReturnCodes(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"TestError", Error, "error"},
		{"TestErrNotOk", ErrNotOk, "not ok status received"},
		{"TestErrFailedMarshalling", ErrFailedMarshalling, "failed to marshal data"},
		{"TestErrFailedUnMarshalling", ErrFailedUnMarshalling, "failed to unmarshal data"},
		{"TestErrProcessingData", ErrProcessingData, "error processing requested data"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.code)
		})
	}
}
