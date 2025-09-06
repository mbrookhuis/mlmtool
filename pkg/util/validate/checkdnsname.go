package validate

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/idna"

	"mlmtool/pkg/util/contains"
	"mlmtool/pkg/util/infoblox"
)

// validates:
//
//	dns name in general
//	if dnsName contains infoblox.DefaultDomain
func CheckDnsName(dnsName string) error {
	// check if dnsName is convertable into valid ascii dns form
	if !isValidFQDN(dnsName) {
		return fmt.Errorf("not a valid fqdn: %s", dnsName)
	}
	if !contains.PartOff(dnsName, infoblox.DefaultDomain) {
		return fmt.Errorf("infoblox record %s not have correct default domain part %s", dnsName, infoblox.DefaultDomain)
	}
	return nil
}

func isValidFQDN(name string) bool {
	// Convert the domain name to ASCII using Punycode
	asciiName, err := idna.ToASCII(name)
	if err != nil {
		return false
	}

	// Check if the ASCII name is within the valid length
	if len(asciiName) > 253 {
		return false
	}

	// Check for consecutive dots or leading/trailing dots
	if strings.Contains(asciiName, "..") || strings.HasPrefix(asciiName, ".") || strings.HasSuffix(asciiName, ".") {
		return false
	}

	// Regular expression for a valid DNS label
	var labelRegex = regexp.MustCompile(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)$`)

	// Split the FQDN into labels and validate each
	labels := strings.Split(asciiName, ".")
	for _, label := range labels {
		if len(label) > 63 || !labelRegex.MatchString(label) {
			return false
		}
	}

	return true
}
