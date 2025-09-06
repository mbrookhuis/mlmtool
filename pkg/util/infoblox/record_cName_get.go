package infoblox

type InfobloxGetRecordCNameRequest struct {
	Name string `json:"name"`
	View string `json:"view"`
}

type InfobloxGetRecordCNameRequestCanonical struct {
	Canonical string `json:"canonical"`
	View      string `json:"view"`
}

type ReplyRecordCName struct {
	Ref       string `json:"_ref"`
	Canonical string `json:"canonical"`
	Name      string `json:"name"`
	View      string `json:"view"`
}

// adGetRecordCName - get cName record
//
// param: name
// param: view
// return:
func (i *Client) adGetRecordCName(name string, view string) ([]ReplyRecordCName, error) {
	request := InfobloxGetRecordCNameRequest{
		Name: name,
		View: view,
	}
	var result []ReplyRecordCName
	err := i.get(cnameRecordAPI, request, &result)
	return result, err
}

// adGetRecordCNameCanonical - get cName record
//
// param: name
// param: view
// return:
func (i *Client) adGetRecordCNameCanonical(canonical string, view string) ([]ReplyRecordCName, error) {
	request := InfobloxGetRecordCNameRequestCanonical{
		Canonical: canonical,
		View:      view,
	}
	var result []ReplyRecordCName
	err := i.get(cnameRecordAPI, request, &result)
	return result, err
}
