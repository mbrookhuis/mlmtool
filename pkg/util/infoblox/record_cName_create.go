package infoblox

// InfobloxCreateCNameReq - Infoblox Create CName req
type InfobloxCreateRecordACNameRequest struct {
	Name      string `json:"name"`
	View      string `json:"view"`
	Canonical string `json:"canonical"`
}

// adCreateRecordCName
//
// param: name
// param: canonical
// param: view
func (i *Client) adCreateRecordCName(name string, canonical string, view string) error {
	request := InfobloxCreateRecordACNameRequest{
		Name:      name,
		Canonical: canonical,
		View:      view,
	}
	return i.post(cnameRecordAPI, request, nil)
}
