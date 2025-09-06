package checkvalue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckEmptyString(t *testing.T) {
	tests := []struct {
		item     string
		val      string
		expected error
	}{
		{"username", "JohnDoe", nil},
		{"email", " ", fmt.Errorf("the value of email is empty")},
		{"password", "", fmt.Errorf("the value of password is empty")},
		{"address", "123 Main St", nil},
		{"phone", "\t\n", fmt.Errorf("the value of phone is empty")},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			err := CheckEmptyString(tt.item, tt.val)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}
