package infoblox

type InfobloxCreateRecordARequest struct {
	Name     string `json:"name"`
	IPv4Addr string `json:"ipv4addr"`
	View     string `json:"view"`
}

// CreateRecordA creates a new A record.
func (i *Client) adCreateRecordA(name string, address string, view string) error {
	request := InfobloxCreateRecordARequest{
		Name:     name,
		IPv4Addr: address,
		View:     view,
	}
	return i.post(aRecordAPI, request, nil)
}
