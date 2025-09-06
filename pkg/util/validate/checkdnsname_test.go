package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"mlmtool/pkg/util/infoblox"
)

// Unit test for CheckDnsName function
func TestCheckDnsName(t *testing.T) {
	cases := []struct {
		name      string
		dnsName   string
		expectErr bool
	}{
		{"Valid DNS Name,Invalid Infoblox domain", "example.com", true},
		{"Valid DNS Name,Valid Infoblox domain", fmt.Sprintf("example.%s", infoblox.DefaultDomain), false},
		{"Invalid DNS Name", ".invalidö.tΩ}@com", true},
		{"Empty DNS Name", "", true},
		{"Infoblox Domain only", infoblox.DefaultDomain, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := CheckDnsName(c.dnsName)
			if c.expectErr {
				assert.Error(t, err) // Ensure error is returned for invalid cases
			} else {
				assert.NoError(t, err) // Ensure no error is returned for valid cases
			}
		})
	}
}
