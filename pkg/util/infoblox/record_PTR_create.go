package infoblox

type InfobloxCreateRecordPTRRequest struct {
	Name     string `json:"name"`
	PtrDName string `json:"ptrdname"`
	IPv4Addr string `json:"ipv4addr"`
	View     string `json:"view"`
}

// CreateRecordPTR creates a new PTR record.
func (i *Client) adCreateRecordPTR(name string, ptrDName string, address string, view string) error {
	request := InfobloxCreateRecordPTRRequest{
		Name:     name,
		PtrDName: ptrDName,
		IPv4Addr: address,
		View:     view,
	}
	return i.post(ptrRecordAPI, request, nil)
}
