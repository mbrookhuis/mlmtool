package contains

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	tests := []struct {
		slice   []string
		element string
		want    bool
	}{
		{[]string{"apple", "banana", "cherry"}, "banana", true},
		{[]string{"apple", "banana", "cherry"}, "grape", false},
		{[]string{}, "banana", false},
	}

	for _, tt := range tests {
		t.Run(tt.element, func(t *testing.T) {
			assert.Equal(t, tt.want, Contains(tt.slice, tt.element))
		})
	}
}

func TestSubInString(t *testing.T) {
	tests := []struct {
		slice  []string
		substr string
		want   bool
	}{
		{[]string{"apple", "banana", "cherry"}, "ban", true},
		{[]string{"apple", "banana", "cherry"}, "berry", false},
		{[]string{}, "ban", false},
	}

	for _, tt := range tests {
		t.Run(tt.substr, func(t *testing.T) {
			assert.Equal(t, tt.want, SubInString(tt.slice, tt.substr))
		})
	}
}

func TestPartOff(t *testing.T) {
	tests := []struct {
		str    string
		substr string
		want   bool
	}{
		{"apple banana cherry", "banana", true},
		{"apple banana cherry", "grape", false},
		{"", "banana", false},
	}

	for _, tt := range tests {
		t.Run(tt.substr, func(t *testing.T) {
			assert.Equal(t, tt.want, PartOff(tt.str, tt.substr))
		})
	}
}

func TestExists(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "exists_test-*.txt")
	tempFilePath := tempFile.Name()
	assert.NoError(t, err)
	defer os.Remove(tempFilePath) // Clean up

	tests := []struct {
		path  string
		valid bool
	}{
		{tempFilePath, true},
		{"/path/to/nonexistent/file", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			exists, _ := Exists(tt.path)
			assert.Equal(t, tt.valid, exists)
		})
	}
}
