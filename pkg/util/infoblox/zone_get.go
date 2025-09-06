package infoblox

type InfobloxGetZoneAuthRequest struct {
	Fqdn string `json:"fqdn"`
	View string `json:"view"`
}

type InfobloxReplyZoneAuth struct {
	Ref  string `json:"_ref"`
	Fqdn string `json:"fqdn"`
	View string `json:"view"`
}

// adGetZoneAuth
//
// param: fqdn
// param: view
// return:
func (i *Client) adGetZoneAuth(fqdn string, view string) ([]InfobloxReplyZoneAuth, error) {
	request := InfobloxGetZoneAuthRequest{
		Fqdn: fqdn,
		View: view,
	}
	var result []InfobloxReplyZoneAuth
	err := i.get(zoneAuthAPI, request, &result)
	return result, err
}
