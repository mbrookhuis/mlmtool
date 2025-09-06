package infoblox

import "fmt"

type InfobloxRestartRequest struct {
	RestartOption   string `json:"restart_option"`
	ServiceOption   string `json:"service_option"`
	MemberOrder     string `json:"member_order"`
	SequentialDelay int    `json:"sequential_delay"`
}

// adRestartInfoblox
//
// param: gridLink
func (i *Client) adRestartInfoblox(gridLink string) error {
	request := InfobloxRestartRequest{
		RestartOption:   "RESTART_IF_NEEDED",
		ServiceOption:   "ALL",
		MemberOrder:     "SEQUENTIALLY",
		SequentialDelay: 1}
	return i.post(fmt.Sprintf("%v%v", gridLink, restartAPI), request, nil)
}
