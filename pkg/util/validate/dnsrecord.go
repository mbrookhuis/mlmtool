package validate

import (
	"fmt"
	"strings"

	"mlmtool/pkg/util/infoblox"
	returncodes "mlmtool/pkg/util/returnCodes"
)

func IBvalidateDNSRecord(p infoblox.DNSRecord, dnsRecordSlice []string) (err error) {
	if err = CheckIP(p.IP); err != nil {
		return fmt.Errorf("%s: %s\n%s", returncodes.NoValidIP, p.IP, err)
	}
	if err = CheckDnsName(p.DNSRecord); err != nil {
		return fmt.Errorf("%s: %s\n%s", returncodes.NoValidDnsName, p.DNSRecord, err)
	}
	for _, dnsName := range p.CName {
		if err = CheckDnsName(dnsName); err != nil {
			return fmt.Errorf("%s: %s\n%s", returncodes.NoValidDnsName, dnsName, err)
		}
	}
	if len(dnsRecordSlice) < len(strings.Split(infoblox.DefaultDomain, "."))+2 {
		return fmt.Errorf("domain seems faulty: %s\nWe expect at least dns record with hostName + DefaultDomain", p.DNSRecord)
	}

	return nil
}
