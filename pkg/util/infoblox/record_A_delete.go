package infoblox

import "fmt"

// adDeleteARecord
//
// param: linkRef
func (i *Client) adDeleteARecord(linkRef string) error {
	err := i.delete(fmt.Sprintf("%v%v", linkRef, aRecordAPI))
	return err
}
