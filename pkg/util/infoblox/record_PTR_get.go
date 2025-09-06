package infoblox

type InfobloxGetRecordPTRRequest struct {
	PtrDName string `json:"ptrdname"`
	View     string `json:"view"`
}

type ReplyRecordPTR struct {
	Ref      string `json:"_ref"`
	PtrDName string `json:"ptrdname"`
	View     string `json:"view"`
}

// adGetRecordPTR
//
// param: name
// param: view
// return:
func (i *Client) adGetRecordPTR(ptrDName string, view string) ([]ReplyRecordPTR, error) {
	request := InfobloxGetRecordPTRRequest{
		PtrDName: ptrDName,
		View:     view,
	}
	var result []ReplyRecordPTR
	err := i.get(ptrRecordAPI, request, &result)
	return result, err
}
