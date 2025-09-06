package infoblox

type InfobloxCreateZoneAuthRequest struct {
	Fqdn string `json:"fqdn"`
	View string `json:"view"`
}

// adCreateZoneAuth
//
// param: domain
// param: view
func (i *Client) adCreateZoneAuth(domain string, view string) error {
	request := InfobloxCreateZoneAuthRequest{
		Fqdn: domain,
		View: view,
	}
	return i.post(zoneAuthAPI, request, nil)
}
