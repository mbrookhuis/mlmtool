package infoblox

type InfobloxGridRequest struct {
	Ref string `json:"_ref"`
}

// adGridInfoblox
func (i *Client) adGridInfoblox() ([]InfobloxGridRequest, error) {
	var result []InfobloxGridRequest
	err := i.get(gridAPI, nil, &result)
	return result, err
}
