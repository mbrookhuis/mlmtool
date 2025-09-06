package validate

import (
	"strings"
	"testing"

	"mlmtool/pkg/util/infoblox"

	"github.com/stretchr/testify/require"
)

func TestIBvalidateDNSRecord(t *testing.T) {
	paramsDnsInvalidCNAME := infoblox.DNSRecord{
		DNSRecord: "host.net.dev." + infoblox.DefaultDomain,
		IP:        "192.168.7.9",
		CName:     []string{"..www.." + "host.net.dev." + infoblox.DefaultDomain},
	}

	paramsDns := infoblox.DNSRecord{
		DNSRecord: "host.dev." + infoblox.DefaultDomain,
		IP:        "192.168.7.9",
		PodRecord: true,
		PtrRecord: true,
		CName: []string{
			"www." + "host.net.dev." + infoblox.DefaultDomain,
			"mail." + "host.net.dev." + infoblox.DefaultDomain,
		},
	}

	paramsInvalidIP := paramsDns
	paramsInvalidIP.IP = "192.168.7..9"

	paramsInvalidDNS := paramsDns
	paramsInvalidDNS.DNSRecord = "test.example.com"

	paramsSuspiciousDNS := paramsDns
	paramsSuspiciousDNS.DNSRecord = "test." + infoblox.DefaultDomain

	tests := []struct {
		name      string
		paramsDns infoblox.DNSRecord
		expectErr bool
	}{
		{
			name:      "Invalid IP",
			paramsDns: paramsInvalidIP,
			expectErr: true,
		},
		{
			name:      "Invalid domain: missing default part",
			paramsDns: paramsInvalidDNS,
			expectErr: true,
		},
		{
			name:      "Suspicious domain length",
			paramsDns: paramsSuspiciousDNS,
			expectErr: true,
		},
		{
			name:      "paramsDnsInvalidCNAME",
			paramsDns: paramsDnsInvalidCNAME,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dnsRecordSlice := strings.Split(tt.paramsDns.DNSRecord, ".")
			err := IBvalidateDNSRecord(tt.paramsDns, dnsRecordSlice)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}
