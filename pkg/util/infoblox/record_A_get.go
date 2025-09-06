package infoblox

type InfobloxGetRecordARequest struct {
	Name string `json:"name"`
	View string `json:"view"`
}

type ReplyRecordA struct {
	Ref      string `json:"_ref"`
	Ipv4Addr string `json:"ipv4addr"`
	Name     string `json:"name"`
	View     string `json:"view"`
}

// adGetRecordA
//
// param: name
// param: view
// return:
func (i *Client) adGetRecordA(name string, view string) ([]ReplyRecordA, error) {
	request := InfobloxGetRecordARequest{
		Name: name,
		View: view,
	}
	var result []ReplyRecordA
	err := i.get(aRecordAPI, request, &result)
	return result, err
}
